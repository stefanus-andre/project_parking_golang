-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Jun 14, 2025 at 10:04 AM
-- Server version: 10.4.32-MariaDB
-- PHP Version: 8.2.12

SET SQL_MODE = "NO_AUTO_VALUE_ON_ZERO";
START TRANSACTION;
SET time_zone = "+00:00";


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;

--
-- Database: `parking_system`
--

-- --------------------------------------------------------

--
-- Table structure for table `parking_history`
--

CREATE TABLE `parking_history` (
  `id` int(11) NOT NULL,
  `registration_no` varchar(50) NOT NULL,
  `slot_number` int(11) NOT NULL,
  `parked_at` timestamp NOT NULL DEFAULT current_timestamp() ON UPDATE current_timestamp(),
  `left_at` timestamp NULL DEFAULT NULL,
  `hours_parked` int(11) DEFAULT 0,
  `charge_amount` decimal(10,2) DEFAULT 0.00
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `parking_history`
--

INSERT INTO `parking_history` (`id`, `registration_no`, `slot_number`, `parked_at`, `left_at`, `hours_parked`, `charge_amount`) VALUES
(1, 'KA-01-HH-1234', 1, '2025-06-14 08:04:29', '2025-06-14 08:04:29', 4, 30.00),
(2, 'KA-01-HH-9999', 2, '2025-06-14 08:04:29', NULL, 0, 0.00),
(3, 'KA-01-BB-0001', 3, '2025-06-14 08:04:29', '2025-06-14 08:04:29', 6, 50.00),
(4, 'KA-01-HH-7777', 4, '2025-06-14 08:04:29', NULL, 0, 0.00),
(5, 'KA-01-HH-2701', 5, '2025-06-14 08:04:29', NULL, 0, 0.00),
(6, 'KA-01-HH-3141', 6, '2025-06-14 08:04:29', '2025-06-14 08:04:29', 4, 30.00),
(7, 'KA-01-P-333', 6, '2025-06-14 08:04:29', NULL, 0, 0.00),
(8, 'KA-09-HH-0987', 1, '2025-06-14 08:04:29', NULL, 0, 0.00),
(9, 'CA-09-IO-1111', 3, '2025-06-14 08:04:29', NULL, 0, 0.00);

-- --------------------------------------------------------

--
-- Table structure for table `parking_lots`
--

CREATE TABLE `parking_lots` (
  `id` int(11) NOT NULL,
  `capacity` int(11) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `parking_lots`
--

INSERT INTO `parking_lots` (`id`, `capacity`, `created_at`) VALUES
(1, 6, '2025-06-14 08:04:29');

-- --------------------------------------------------------

--
-- Table structure for table `parking_slots`
--

CREATE TABLE `parking_slots` (
  `slot_number` int(11) NOT NULL,
  `registration_no` varchar(50) DEFAULT NULL,
  `is_occupied` tinyint(1) DEFAULT 0,
  `parked_at` timestamp NULL DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `parking_slots`
--

INSERT INTO `parking_slots` (`slot_number`, `registration_no`, `is_occupied`, `parked_at`) VALUES
(1, 'KA-09-HH-0987', 1, '2025-06-14 08:04:29'),
(2, 'KA-01-HH-9999', 1, '2025-06-14 08:04:29'),
(3, 'CA-09-IO-1111', 1, '2025-06-14 08:04:29'),
(4, 'KA-01-HH-7777', 1, '2025-06-14 08:04:29'),
(5, 'KA-01-HH-2701', 1, '2025-06-14 08:04:29'),
(6, 'KA-01-P-333', 1, '2025-06-14 08:04:29');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `parking_history`
--
ALTER TABLE `parking_history`
  ADD PRIMARY KEY (`id`),
  ADD KEY `idx_registration_history` (`registration_no`);

--
-- Indexes for table `parking_lots`
--
ALTER TABLE `parking_lots`
  ADD PRIMARY KEY (`id`);

--
-- Indexes for table `parking_slots`
--
ALTER TABLE `parking_slots`
  ADD PRIMARY KEY (`slot_number`),
  ADD KEY `idx_registration` (`registration_no`),
  ADD KEY `idx_occupied` (`is_occupied`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `parking_history`
--
ALTER TABLE `parking_history`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=10;

--
-- AUTO_INCREMENT for table `parking_lots`
--
ALTER TABLE `parking_lots`
  MODIFY `id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
