-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Mar 12, 2025 at 06:48 AM
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
-- Database: `mental_health_v2`
--

-- --------------------------------------------------------

--
-- Table structure for table `konsultan`
--

CREATE TABLE `konsultan` (
  `konsultan_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `pengalaman` text NOT NULL,
  `spesialisasi` varchar(100) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `pesan_konsultasi`
--

CREATE TABLE `pesan_konsultasi` (
  `pesan_id` int(11) NOT NULL,
  `sesi_id` int(11) NOT NULL,
  `isi_pesan` text NOT NULL,
  `timestamp` datetime DEFAULT current_timestamp(),
  `sentiment_score` float DEFAULT NULL CHECK (`sentiment_score` between -1 and 1),
  `urgency_level` enum('rendah','sedang','tinggi') NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `respons_chatbot`
--

CREATE TABLE `respons_chatbot` (
  `respons_id` int(11) NOT NULL,
  `pesan_id` int(11) NOT NULL,
  `isi_respons` text NOT NULL,
  `timestamp` datetime DEFAULT current_timestamp()
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `sesi_konsultasi`
--

CREATE TABLE `sesi_konsultasi` (
  `sesi_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `konsultan_id` int(11) DEFAULT NULL,
  `tipe_konsultasi` enum('Chatbot','Konsultan') NOT NULL,
  `status` enum('Berlangsung','Selesai','Dibatalkan') NOT NULL DEFAULT 'Berlangsung',
  `waktu_mulai` datetime DEFAULT current_timestamp(),
  `waktu_selesai` datetime DEFAULT NULL,
  `auto_forward` tinyint(1) DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

-- --------------------------------------------------------

--
-- Table structure for table `user`
--

CREATE TABLE `user` (
  `user_id` int(11) NOT NULL,
  `nama` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password` varchar(255) NOT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `user`
--

INSERT INTO `user` (`user_id`, `nama`, `email`, `password`) VALUES
(1, 'agus', 'agus@gmail.com', 'agus');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `konsultan`
--
ALTER TABLE `konsultan`
  ADD PRIMARY KEY (`konsultan_id`),
  ADD UNIQUE KEY `user_id` (`user_id`);

--
-- Indexes for table `pesan_konsultasi`
--
ALTER TABLE `pesan_konsultasi`
  ADD PRIMARY KEY (`pesan_id`),
  ADD KEY `sesi_id` (`sesi_id`);

--
-- Indexes for table `respons_chatbot`
--
ALTER TABLE `respons_chatbot`
  ADD PRIMARY KEY (`respons_id`),
  ADD KEY `pesan_id` (`pesan_id`);

--
-- Indexes for table `sesi_konsultasi`
--
ALTER TABLE `sesi_konsultasi`
  ADD PRIMARY KEY (`sesi_id`),
  ADD KEY `user_id` (`user_id`),
  ADD KEY `konsultan_id` (`konsultan_id`);

--
-- Indexes for table `user`
--
ALTER TABLE `user`
  ADD PRIMARY KEY (`user_id`),
  ADD UNIQUE KEY `email` (`email`);

--
-- AUTO_INCREMENT for dumped tables
--

--
-- AUTO_INCREMENT for table `konsultan`
--
ALTER TABLE `konsultan`
  MODIFY `konsultan_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `pesan_konsultasi`
--
ALTER TABLE `pesan_konsultasi`
  MODIFY `pesan_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `respons_chatbot`
--
ALTER TABLE `respons_chatbot`
  MODIFY `respons_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `sesi_konsultasi`
--
ALTER TABLE `sesi_konsultasi`
  MODIFY `sesi_id` int(11) NOT NULL AUTO_INCREMENT;

--
-- AUTO_INCREMENT for table `user`
--
ALTER TABLE `user`
  MODIFY `user_id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- Constraints for dumped tables
--

--
-- Constraints for table `konsultan`
--
ALTER TABLE `konsultan`
  ADD CONSTRAINT `konsultan_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE CASCADE;

--
-- Constraints for table `pesan_konsultasi`
--
ALTER TABLE `pesan_konsultasi`
  ADD CONSTRAINT `pesan_konsultasi_ibfk_1` FOREIGN KEY (`sesi_id`) REFERENCES `sesi_konsultasi` (`sesi_id`) ON DELETE CASCADE;

--
-- Constraints for table `respons_chatbot`
--
ALTER TABLE `respons_chatbot`
  ADD CONSTRAINT `respons_chatbot_ibfk_1` FOREIGN KEY (`pesan_id`) REFERENCES `pesan_konsultasi` (`pesan_id`) ON DELETE CASCADE;

--
-- Constraints for table `sesi_konsultasi`
--
ALTER TABLE `sesi_konsultasi`
  ADD CONSTRAINT `sesi_konsultasi_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE CASCADE,
  ADD CONSTRAINT `sesi_konsultasi_ibfk_2` FOREIGN KEY (`konsultan_id`) REFERENCES `konsultan` (`konsultan_id`) ON DELETE SET NULL;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
