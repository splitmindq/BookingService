CREATE SCHEMA IF NOT EXISTS booking_service;

SET search_path TO booking_service;

CREATE TYPE role_type AS ENUM ('admin');
CREATE TYPE booking_status AS ENUM ('pending', 'confirmed', 'canceled');
CREATE TYPE payment_method AS ENUM ('card', 'bank');
CREATE TYPE payment_status AS ENUM ('pending', 'success', 'failed');
CREATE TYPE blocked_reason AS ENUM ('booking', 'maintenance', 'personal', 'other');

-- Users Table
CREATE TABLE IF NOT EXISTS booking_service.users (
    user_id SERIAL PRIMARY KEY,
    phone VARCHAR(20) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Roles Table
CREATE TABLE IF NOT EXISTS booking_service.user_roles(
    user_id INTEGER NOT NULL,
    role role_type NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, role),
    FOREIGN KEY (user_id) REFERENCES booking_service.users (user_id)
);

-- Properties Table
CREATE TABLE IF NOT EXISTS booking_service.properties (
    property_id SERIAL PRIMARY KEY,
    owner_id INTEGER NOT NULL,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    address VARCHAR(255) NOT NULL,
    city VARCHAR(50) NOT NULL,
    country VARCHAR(50) NOT NULL,
    price_per_night NUMERIC(10,2) NOT NULL,
    max_guests INTEGER NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (owner_id) REFERENCES booking_service.users (user_id)
    );

-- Bookings Table
CREATE TABLE IF NOT EXISTS booking_service.bookings (
    booking_id SERIAL PRIMARY KEY,
    property_id INTEGER NOT NULL,
    guest_id INTEGER NOT NULL,
    check_in_date DATE NOT NULL,
    check_out_date DATE NOT NULL,
    total_price NUMERIC(10,2) NOT NULL,
    status booking_status NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (property_id) REFERENCES booking_service.properties (property_id),
    FOREIGN KEY (guest_id) REFERENCES booking_service.users (user_id)
    );

CREATE INDEX idx_bookings_dates ON booking_service.bookings (property_id, check_in_date, check_out_date);

-- Payments Table
CREATE TABLE IF NOT EXISTS booking_service.payments (
    payment_id SERIAL PRIMARY KEY,
    booking_id INTEGER NOT NULL,
    amount NUMERIC(10,2) NOT NULL,
    currency CHAR(3) NOT NULL DEFAULT 'RUB',
    payment_method payment_method NOT NULL,
    payment_status payment_status NOT NULL DEFAULT 'pending',
    transaction_id VARCHAR(100),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (booking_id) REFERENCES booking_service.bookings (booking_id)
    );

-- Table Reviews 
CREATE TABLE IF NOT EXISTS booking_service.reviews (
    review_id SERIAL PRIMARY KEY,
    property_id INTEGER NOT NULL,
    guest_id INTEGER NOT NULL,
    rating INTEGER NOT NULL CHECK (rating BETWEEN 1 AND 5),
    comment TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (property_id) REFERENCES booking_service.properties (property_id),
    FOREIGN KEY (guest_id) REFERENCES booking_service.users (user_id)
    );

-- Table BlockedDates
CREATE TABLE IF NOT EXISTS booking_service.blocked_dates (
    blocked_date_id SERIAL PRIMARY KEY,
    property_id INTEGER NOT NULL,
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    reason blocked_reason NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (property_id) REFERENCES booking_service.properties (property_id)
    );

CREATE INDEX idx_blocked_dates ON booking_service.blocked_dates (property_id, start_date, end_date);