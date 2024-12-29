CREATE TABLE `users` (
    `id` int UNSIGNED NOT NULL,
    `email` varchar(255) NOT NULL,
    `password` char(60) NOT NULL,
    PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `invoices` (
    `id` int UNSIGNED NOT NULL,
    `user_id` int UNSIGNED NOT NULL,
    `bill_to_first_name` varchar(255) NOT NULL,
    `bill_to_last_name` varchar(255) NOT NULL,
    `pay_to_first_name` varchar(255) NOT NULL,
    `pay_to_last_name` varchar(255) NOT NULL,
    `amount_due` int UNSIGNED NOT NULL,
    `amount_paid` int UNSIGNED NOT NULL,
    `status` enum('pending') NOT NULL,
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