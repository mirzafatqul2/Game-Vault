package user

import (
	"Mini-Project-Game-Vault-API/repository/mailjet"
	"Mini-Project-Game-Vault-API/service/wallet"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"github.com/pobyzaarif/goshortcute"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	logger                  *slog.Logger
	userRepo                UserRepository
	appDeploymentUrl        string
	jwtSign                 string
	appEmailVerificationKey string
	notifRepo               *mailjet.MailjetRepository
	walletRepo              wallet.WalletRepository
}

const (
	VerificationCodeTTL = 5
)

type UserService interface {
	Register(ctx context.Context, user User) (string, error)
	VerifyEmail(ctx context.Context, verificationCodeEncrypt string) error
	Login(ctx context.Context, email string, password string) (string, error)
	GetUserProfile(ctx context.Context, userID string) (User, error)
	DepositAmount(ctx context.Context, userID string, description string, amount int) (User, error)
}

func NewUserService(
	logger *slog.Logger,
	userRepo UserRepository,
	appDeploymentUrl string,
	jwtSign string,
	appEmailVerificationKey string,
	notifRepo *mailjet.MailjetRepository,
	walletRepo wallet.WalletRepository,
) UserService {
	return &userService{
		logger:                  logger,
		userRepo:                userRepo,
		appDeploymentUrl:        appDeploymentUrl,
		jwtSign:                 jwtSign,
		appEmailVerificationKey: appEmailVerificationKey,
		notifRepo:               notifRepo,
		walletRepo:              walletRepo,
	}
}

const (
	SubjectRegisterAccount   = "Activate Your Account!"
	EmailBodyRegisterAccount = `Halo, %v, Aktivasi akun anda dengan membuka tautan dibawah<br/><br/>%v<br/>catatan: link hanya berlaku %v menit`
)

func (s *userService) Register(ctx context.Context, user User) (string, error) {
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))
	user.Username = strings.TrimSpace(user.Username)
	// Check email
	if existing, err := s.userRepo.GetByEmail(ctx, user.Email); err == nil && existing.Email != "" {
		return "", errors.New("email already registered")
	}

	// Check username
	if existing, err := s.userRepo.GetByUsername(ctx, user.Username); err == nil && existing.Username != "" {
		return "", errors.New("username already registered")
	}

	// Hashing password
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	user.ID = uuid.New().String()
	user.Password = string(hashPassword)
	user.Role = "user"

	if err = s.userRepo.Create(ctx, user); err != nil {
		return "", err
	}

	timeNow := time.Now()
	expAt := timeNow.Add(time.Duration(time.Minute * VerificationCodeTTL)).Unix()

	verificationCode := fmt.Sprintf("%v|%v", user.Email, expAt)
	verificationCodeEncrypt, _ := goshortcute.AESCBCEncrypt([]byte(verificationCode), []byte(s.appEmailVerificationKey))
	verifCode := goshortcute.StringtoBase64Encode(verificationCodeEncrypt)
	activationLink := s.appDeploymentUrl + "/users/email-verification/" + verifCode

	err = s.notifRepo.SendEmail(user.FullName, user.Email, SubjectRegisterAccount, fmt.Sprintf(EmailBodyRegisterAccount, user.FullName, activationLink, VerificationCodeTTL))
	if err != nil {
		s.logger.Error("send email failed", slog.Any("err", err))
	}

	return user.ID, nil
}

func (s *userService) VerifyEmail(ctx context.Context, verificationCodeEncrypt string) error {
	verifCodeDecode := goshortcute.StringtoBase64Decode(verificationCodeEncrypt)
	verificationCodeDecrypt, err := goshortcute.AESCBCDecrypt([]byte(verifCodeDecode), []byte(s.appEmailVerificationKey))
	if err != nil {
		s.logger.Error("verify email err", slog.Any("err", err.Error()))
		return errors.New("invalid or expired url")
	}

	verificationCode := strings.Split(verificationCodeDecrypt, "|")
	if len(verificationCode) != 2 {
		s.logger.Error("verify email err", slog.Any("err", verificationCodeDecrypt))
		return errors.New("invalid or expired url")
	}

	email := verificationCode[0]
	expAtStr := verificationCode[1]

	ts, err := strconv.ParseInt(expAtStr, 10, 64)
	if err != nil {
		s.logger.Error("verify email err", slog.Any("err", verificationCodeDecrypt))
		return errors.New("invalid or expired url")
	}
	expAt := time.Unix(ts, 0)
	if time.Now().After(expAt) {
		return errors.New("invalid or expired url")
	}

	getUser, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		s.logger.Error("verify email err", slog.Any("err", err))
		return err
	}

	if getUser.IsVerified {
		s.logger.Warn("verify email err", slog.Any("err", "email verified already"))
		return errors.New("invalid or expired url")
	}

	getUser.IsVerified = true
	if err := s.userRepo.UpdatEmailVerification(ctx, getUser); err != nil {
		s.logger.Error("verify email err", slog.Any("err", err))
		return err
	}

	return nil
}

func (s *userService) Login(ctx context.Context, email string, password string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil || user.Email == "" {
		return "", errors.New("invalid email or password")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		s.logger.Error("login err", slog.Any("err", err.Error()))
		return "", errors.New("invalid email or password")
	}

	if !user.IsVerified {
		return "", errors.New("email address has not been verified")
	}

	token, err := s.generateToken(s.jwtSign, user.ID, user.Role)
	if err != nil {
		s.logger.Error("generate token err", slog.Any("err", err.Error()))
		return "", errors.New("generate token error")
	}

	return token, err
}

func (s *userService) generateToken(jwtSign string, id string, role string) (string, error) {
	type jwtClaims struct {
		ID   string `json:"id"`
		Role string `json:"role"`
		jwt.RegisteredClaims
	}

	timeNow := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims{
		ID:   id,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(timeNow),
			ExpiresAt: jwt.NewNumericDate(timeNow.Add(time.Hour * 24)),
		},
	})

	signedToken, err := token.SignedString([]byte(jwtSign))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (s *userService) GetUserProfile(ctx context.Context, userID string) (User, error) {
	return s.userRepo.GetByID(ctx, userID)
}

func (s *userService) DepositAmount(ctx context.Context, userID, description string, amount int) (User, error) {
	if amount <= 0 {
		return User{}, errors.New("invalid deposit amount")
	}

	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return User{}, err
	}

	newBalance := user.DepositAmount + amount

	err = s.userRepo.UpdateDeposit(ctx, userID, newBalance)
	if err != nil {
		return User{}, err
	}

	wallet := wallet.WalletTransaction{
		ID:          uuid.New().String(),
		UserID:      userID,
		Type:        "deposit",
		Amount:      amount,
		Description: description,
	}

	s.walletRepo.Create(ctx, wallet)

	return user, err
}
