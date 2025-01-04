CREATE TABLE `users` (
    `id` int UNSIGNED NOT NULL,
    `email` varchar(255) NOT NULL,
    `password` char(60) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `invoices` (
    `id` int UNSIGNED NOT NULL,
    `user_id` int UNSIGNED NOT NULL,
    `invoice_number` varchar(50) NOT NULL,
    `po_number` varchar(50) NOT NULL,
    `currency` char(3) NOT NULL,
    `due_date` date NOT NULL,
    `message` varchar(255) NOT NULL,
    `bill_to_first_name` varchar(255) NOT NULL,
    `bill_to_last_name` varchar(255) NOT NULL,
    `bill_to_company` varchar(255) NOT NULL,
    `bill_to_address_line_1` varchar(255) NOT NULL,
    `bill_to_address_line_2` varchar(255) NOT NULL,
    `bill_to_city` varchar(255) NOT NULL,
    `bill_to_state` varchar(3) NOT NULL,
    `bill_to_postal_code` varchar(12) NOT NULL,
    `bill_to_country` char(2) NOT NULL,
    `bill_to_email` varchar(255) NOT NULL,
    `bill_to_phone` varchar(15) NOT NULL,
    `pay_to_first_name` varchar(255) NOT NULL,
    `pay_to_last_name` varchar(255) NOT NULL,
    `pay_to_company` varchar(255) NOT NULL,
    `pay_to_address_line_1` varchar(255) NOT NULL,
    `pay_to_address_line_2` varchar(255) NOT NULL,
    `pay_to_city` varchar(255) NOT NULL,
    `pay_to_state` varchar(3) NOT NULL,
    `pay_to_postal_code` varchar(12) NOT NULL,
    `pay_to_country` char(2) NOT NULL,
    `pay_to_email` varchar(255) NOT NULL,
    `pay_to_phone` varchar(15) NOT NULL,
    `line_items` json DEFAULT NULL,
    `payment_methods` json DEFAULT NULL,
    `tax_rate` varchar(10) NOT NULL,
    `amount_due` int UNSIGNED NOT NULL,
    `amount_paid` int UNSIGNED NOT NULL,
    `status` enum('pending') NOT NULL,
    `created_at` datetime NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `transactions` (
    `id` int UNSIGNED NOT NULL,
    `user_id` int UNSIGNED NOT NULL,
    `type` enum('authorize', 'capture', 'sale', 'void', 'refund') NOT NULL,
    `card_type` varchar(255) NOT NULL,
    `amount_captured` int UNSIGNED NOT NULL,
    `invoice_id` int UNSIGNED NOT NULL,
    `status` enum('approved', 'declined') NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;