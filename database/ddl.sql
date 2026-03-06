-- USERS
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    full_name VARCHAR(150) NOT NULL,
    username VARCHAR(150) UNIQUE NOT NULL,
    email VARCHAR(150) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    deposit_amount INT NOT NULL DEFAULT 0,
    role VARCHAR(20) NOT NULL DEFAULT 'user',
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE email_verifications (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- GAMES
CREATE TABLE games (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    external_id VARCHAR(100) UNIQUE,
    name VARCHAR(200) NOT NULL,
    category VARCHAR(100) NOT NULL,
    stock_availability INT NOT NULL DEFAULT 0,
    rental_cost INT NOT NULL,
    image_url TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- RENTALS
CREATE TABLE rentals (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    game_id UUID NOT NULL,
    rental_days INT NOT NULL,
    total_cost INT NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'ongoing',
    rented_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    due_date TIMESTAMP NOT NULL,
    returned_at TIMESTAMP,
    penalty_fee INT DEFAULT 0,

    CONSTRAINT fk_rental_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_rental_game
        FOREIGN KEY(game_id)
        REFERENCES games(id)
        ON DELETE CASCADE
);

-- WALLETS
CREATE TABLE wallet_transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    amount INT NOT NULL,
    type VARCHAR(20) NOT NULL,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_wallet_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE

);
    
-- WHISTLISTS
CREATE TABLE wishlists (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    game_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_wishlist_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_wishlist_game
        FOREIGN KEY(game_id)
        REFERENCES games(id)
        ON DELETE CASCADE,

    CONSTRAINT unique_user_game_wishlist
        UNIQUE(user_id, game_id)
);

-- REVIEWS
CREATE TABLE reviews (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL,
    game_id UUID NOT NULL,
    rating INT CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT fk_review_user
        FOREIGN KEY(user_id)
        REFERENCES users(id)
        ON DELETE CASCADE,

    CONSTRAINT fk_review_game
        FOREIGN KEY(game_id)
        REFERENCES games(id)
        ON DELETE CASCADE,

    CONSTRAINT unique_user_game_review
        UNIQUE(user_id, game_id)
);