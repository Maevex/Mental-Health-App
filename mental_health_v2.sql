-- phpMyAdmin SQL Dump
-- version 5.2.1
-- https://www.phpmyadmin.net/
--
-- Host: 127.0.0.1
-- Generation Time: Apr 19, 2025 at 02:01 AM
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
-- Table structure for table `konsultan_kontak`
--

CREATE TABLE `konsultan_kontak` (
  `konsultan_id` int(11) NOT NULL,
  `nama` varchar(100) NOT NULL,
  `spesialisasi` varchar(255) DEFAULT NULL,
  `pengalaman` text DEFAULT NULL,
  `no_telepon` varchar(20) DEFAULT NULL,
  `email` varchar(100) DEFAULT NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `konsultan_kontak`
--

INSERT INTO `konsultan_kontak` (`konsultan_id`, `nama`, `spesialisasi`, `pengalaman`, `no_telepon`, `email`) VALUES
(1, 'Dr. Anisa Sari', 'Psikologi Klinis', '10 tahun', '081234567890', 'anisa.sari@example.com');

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

--
-- Dumping data for table `pesan_konsultasi`
--

INSERT INTO `pesan_konsultasi` (`pesan_id`, `sesi_id`, `isi_pesan`, `timestamp`, `sentiment_score`, `urgency_level`) VALUES
(1, 4, 'im living a happy life.', '2025-03-30 08:38:00', 0.571885, 'rendah'),
(2, 4, 'i want to kill myself.', '2025-03-30 08:39:27', -0.659728, 'rendah'),
(3, 4, 'i live a happy life.', '2025-03-30 08:53:48', 0.571885, 'rendah'),
(4, 4, 'i live a happy life.', '2025-03-30 08:58:43', 0.571885, 'rendah'),
(5, 4, 'i live a happy life.', '2025-03-30 09:06:22', 0.571885, 'rendah'),
(6, 4, 'i live a happy life, i love my family.', '2025-03-30 09:16:53', 0.835976, 'rendah'),
(7, 5, 'Halo, ini chat pertama!', '2025-04-18 13:45:49', 0, 'rendah'),
(8, 6, 'Halo, aku ilham', '2025-04-18 13:46:13', 0, 'rendah'),
(9, 4, 'Halo, aku ilham', '2025-04-18 13:46:47', 0, 'rendah'),
(10, 5, 'apa itu mental health?', '2025-04-18 13:47:49', 0, 'rendah'),
(11, 7, 'apa itu mental health?', '2025-04-18 13:54:56', 0, 'rendah'),
(12, 7, 'bagaimana cara menjaga mental health?', '2025-04-18 14:05:37', 0, 'rendah'),
(13, 7, 'bagaimana cara menjaga mental health?', '2025-04-18 17:28:15', 0, 'rendah'),
(14, 7, 'bagaimana cara menjaga mental health?', '2025-04-18 17:31:29', 0, 'rendah');

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

--
-- Dumping data for table `respons_chatbot`
--

INSERT INTO `respons_chatbot` (`respons_id`, `pesan_id`, `isi_respons`, `timestamp`) VALUES
(1, 1, 'Tidak ada respons dari chatbot.', '2025-03-30 08:38:00'),
(2, 3, 'Tidak ada respons dari chatbot.', '2025-03-30 08:53:50'),
(3, 4, 'Tidak ada respons dari chatbot.', '2025-03-30 08:58:43'),
(4, 5, 'That\'s wonderful to hear!  I\'m glad you\'re happy.  Is there anything you\'d like to talk about or ask me about today?\n', '2025-03-30 09:06:22'),
(5, 6, 'That\'s wonderful to hear!  Cherish those relationships and enjoy your happy life.\n', '2025-03-30 09:16:53'),
(6, 7, 'Halo! Senang bertemu denganmu dalam chat pertama ini. Ada yang bisa kubantu?\n', '2025-04-18 13:45:49'),
(7, 8, 'Halo Ilham! Ada yang bisa saya bantu?\n', '2025-04-18 13:46:13'),
(8, 9, 'Halo Ilham! Ada yang bisa saya bantu?\n', '2025-04-18 13:46:47'),
(9, 10, 'Kesehatan mental mengacu pada kesejahteraan emosional, psikologis, dan sosial seseorang. Ini memengaruhi cara kita berpikir, merasa, dan bertindak. Ini juga menentukan bagaimana kita menangani stres, berhubungan dengan orang lain, dan membuat pilihan.  Kesehatan mental penting di setiap tahap kehidupan, dari masa kanak-kanak dan remaja hingga dewasa.\n\nBerikut beberapa aspek kunci dari kesehatan mental yang baik:\n\n* **Kesejahteraan emosional:**  Merasa baik tentang diri sendiri dan mampu mengelola emosi secara sehat, termasuk merasakan dan mengekspresikan berbagai emosi seperti kegembiraan, kesedihan, dan kemarahan dengan cara yang tepat.\n* **Kesejahteraan psikologis:**  Memiliki harga diri yang sehat, merasa kompeten, dan mampu mengatasi tantangan hidup. Ini juga termasuk kemampuan untuk berpikir jernih, berkonsentrasi, dan membuat keputusan.\n* **Kesejahteraan sosial:**  Memiliki hubungan yang sehat dan positif dengan orang lain, merasa terhubung dengan komunitas, dan mampu memberikan kontribusi kepada masyarakat.\n\nPenting untuk diingat bahwa kesehatan mental bukan hanya tentang tidak adanya penyakit mental. Seperti kesehatan fisik, kesehatan mental adalah sebuah *kontinum*.  Seseorang dapat memiliki kesehatan mental yang baik, mengalami masalah kesehatan mental ringan, atau menderita penyakit mental yang serius.\n\nGangguan kesehatan mental adalah kondisi yang memengaruhi pemikiran, perasaan, suasana hati, atau perilaku seseorang.  Beberapa contoh gangguan kesehatan mental meliputi depresi, kecemasan, gangguan bipolar, skizofrenia, dan gangguan stres pascatrauma (PTSD).\n\nSama seperti kita merawat kesehatan fisik kita, penting juga untuk merawat kesehatan mental kita.  Ada banyak hal yang dapat kita lakukan untuk meningkatkan dan menjaga kesehatan mental kita, termasuk:\n\n* **Berolahraga secara teratur:** Aktivitas fisik telah terbukti meningkatkan suasana hati dan mengurangi stres.\n* **Makan makanan yang sehat:** Diet seimbang dapat membantu menjaga kesehatan fisik dan mental.\n* **Tidur yang cukup:** Kurang tidur dapat berdampak negatif pada kesehatan mental.\n* **Menghindari alkohol dan obat-obatan:** Penyalahgunaan zat dapat memperburuk masalah kesehatan mental.\n* **Berhubungan dengan orang lain:**  Memiliki jaringan dukungan sosial yang kuat sangat penting untuk kesehatan mental yang baik.\n* **Melakukan aktivitas yang menyenangkan:**  Meluangkan waktu untuk hobi dan aktivitas yang kita nikmati dapat membantu mengurangi stres dan meningkatkan suasana hati.\n* **Mencari bantuan profesional jika diperlukan:**  Jangan ragu untuk mencari bantuan dari profesional kesehatan mental jika Anda mengalami masalah kesehatan mental.\n\n\nJika Anda atau seseorang yang Anda kenal sedang berjuang dengan masalah kesehatan mental, penting untuk mencari bantuan.  Ada banyak sumber daya yang tersedia, termasuk profesional kesehatan mental, hotline krisis, dan kelompok pendukung.\n', '2025-04-18 13:47:49'),
(10, 11, 'Kesehatan mental mengacu pada kesejahteraan emosional, psikologis, dan sosial seseorang. Ini memengaruhi cara kita berpikir, merasa, dan bertindak. Ini juga membantu menentukan bagaimana kita menangani stres, berhubungan dengan orang lain, dan membuat pilihan. Kesehatan mental penting di setiap tahap kehidupan, mulai dari masa kanak-kanak dan remaja hingga dewasa.\n\nKesehatan mental bukanlah sekadar tidak adanya penyakit mental.  Ini adalah keadaan sejahtera di mana seseorang menyadari kemampuannya sendiri, dapat mengatasi tekanan normal dalam hidup, dapat bekerja secara produktif dan bermanfaat, dan dapat memberikan kontribusi kepada komunitasnya.\n\nBerikut beberapa aspek kunci dari kesehatan mental yang baik:\n\n* **Kesejahteraan emosional:** Merasa baik tentang diri sendiri dan orang lain, mampu mengelola emosi yang sulit, dan memiliki rasa tujuan hidup.\n* **Kesejahteraan psikologis:**  Memiliki pemikiran yang jernih, mampu berkonsentrasi dan belajar, serta memiliki harga diri yang sehat.\n* **Kesejahteraan sosial:**  Memiliki hubungan yang sehat dengan orang lain, mampu membangun dan mempertahankan hubungan, serta merasa terhubung dengan komunitas.\n\nSama seperti kesehatan fisik, kesehatan mental juga dapat terganggu.  Faktor-faktor seperti genetika, pengalaman hidup, dan lingkungan sekitar dapat memengaruhi kesehatan mental seseorang.  Gangguan kesehatan mental, seperti depresi, kecemasan, dan skizofrenia, adalah kondisi medis nyata yang memerlukan perawatan profesional.\n\nMerawat kesehatan mental sama pentingnya dengan merawat kesehatan fisik.  Ada banyak hal yang dapat dilakukan untuk menjaga kesehatan mental, termasuk:\n\n* **Berolahraga secara teratur.**\n* **Makan makanan yang sehat.**\n* **Tidur yang cukup.**\n* **Menghabiskan waktu bersama orang-orang yang dicintai.**\n* **Melakukan aktivitas yang menyenangkan.**\n* **Belajar mengelola stres.**\n* **Mencari bantuan profesional jika diperlukan.**\n\n\nJika Anda merasa kesulitan dengan kesehatan mental Anda, jangan ragu untuk mencari bantuan.  Berbicara dengan seorang profesional kesehatan mental dapat membantu Anda mengembangkan strategi koping dan mendapatkan kembali kesejahteraan Anda.\n', '2025-04-18 13:54:56'),
(11, 12, 'Menjaga kesehatan mental sama pentingnya dengan menjaga kesehatan fisik. Berikut beberapa cara menjaga kesehatan mental:\n\n**1. Pola Hidup Sehat:**\n\n* **Tidur yang Cukup:** Usahakan tidur 7-9 jam setiap malam. Kualitas tidur yang baik sangat penting untuk kesehatan mental.\n* **Makan Makanan Bergizi:** Konsumsi makanan sehat dan seimbang, batasi makanan olahan, gula, dan kafein.\n* **Olahraga Teratur:** Aktivitas fisik melepaskan endorfin yang dapat meningkatkan mood dan mengurangi stres.  Cukup 30 menit sehari beberapa kali seminggu.\n* **Hindari Alkohol dan Narkoba:**  Zat-zat ini dapat memperburuk masalah kesehatan mental.\n* **Paparan Sinar Matahari:** Sinar matahari membantu tubuh memproduksi vitamin D, yang berperan penting dalam kesehatan mental.\n\n**2. Manajemen Stres:**\n\n* **Identifikasi Pemicu Stres:** Kenali hal-hal yang membuat Anda stres dan cari cara untuk menghindarinya atau mengelolanya.\n* **Teknik Relaksasi:**  Latihan pernapasan, meditasi, yoga, dan mindfulness dapat membantu mengurangi stres.\n* **Manajemen Waktu yang Efektif:**  Buat jadwal dan prioritaskan tugas untuk mengurangi rasa kewalahan.\n* **Istirahat yang Cukup:**  Sediakan waktu untuk bersantai dan melakukan aktivitas yang Anda nikmati.\n\n**3. Hubungan Sosial yang Positif:**\n\n* **Jalin Hubungan yang Sehat:**  Habiskan waktu dengan orang-orang yang mendukung dan peduli.\n* **Komunikasi yang Efektif:**  Ekspresikan perasaan Anda dengan jujur dan terbuka.\n* **Bergabung dengan Komunitas:**  Cari komunitas yang sesuai dengan minat Anda untuk memperluas jaringan sosial.\n* **Bantu Orang Lain:**  Membantu orang lain dapat meningkatkan rasa bahagia dan kepuasan diri.\n\n**4. Pengembangan Diri:**\n\n* **Tetapkan Tujuan yang Realistis:**  Memiliki tujuan memberi Anda arah dan motivasi.\n* **Pelajari Keterampilan Baru:**  Mengembangkan keterampilan baru dapat meningkatkan rasa percaya diri.\n* **Lakukan Hobi:**  Luangkan waktu untuk melakukan hal-hal yang Anda sukai.\n* **Bersyukur:**  Fokus pada hal-hal positif dalam hidup Anda.\n\n**5. Mencari Bantuan Profesional:**\n\n* **Jangan Ragu untuk Mencari Bantuan:**  Jika Anda merasa kesulitan mengatasi masalah kesehatan mental, jangan ragu untuk mencari bantuan profesional.\n* **Terapi:**  Terapi dapat membantu Anda memahami dan mengatasi masalah kesehatan mental.\n* **Psikiater:**  Psikiater dapat mendiagnosis dan meresepkan obat-obatan jika diperlukan.\n\n\n**Tanda-tanda Perlu Mencari Bantuan Profesional:**\n\n* Perubahan mood yang drastis dan berkepanjangan.\n* Kesulitan berkonsentrasi.\n* Kehilangan minat pada aktivitas yang biasanya dinikmati.\n* Gangguan tidur atau pola makan.\n* Perasaan cemas atau takut yang berlebihan.\n* Pikiran untuk menyakiti diri sendiri atau orang lain.\n\n\nMenjaga kesehatan mental adalah proses yang berkelanjutan.  Lakukan langkah-langkah di atas secara konsisten dan jangan ragu untuk mencari bantuan jika diperlukan.  Ingat, Anda tidak sendirian.\n', '2025-04-18 14:05:37'),
(12, 13, 'Menjaga kesehatan mental sama pentingnya dengan menjaga kesehatan fisik.  Berikut beberapa cara untuk menjaga kesehatan mental:\n\n**1. Pola Hidup Sehat:**\n\n* **Tidur yang cukup:** Usahakan tidur 7-9 jam setiap malam.  Kualitas tidur yang baik sangat berpengaruh pada suasana hati dan kemampuan kognitif.\n* **Makan makanan bergizi:** Konsumsi makanan sehat dan seimbang, batasi asupan gula, kafein, dan makanan olahan.\n* **Olahraga teratur:** Aktivitas fisik melepaskan endorfin yang dapat meningkatkan mood dan mengurangi stres.  Cukup 30 menit olahraga ringan setiap hari atau beberapa kali seminggu.\n* **Hindari alkohol dan narkoba:**  Substansi ini dapat memperburuk masalah kesehatan mental yang sudah ada atau memicu masalah baru.\n* **Terhidrasi dengan baik:** Minum cukup air putih setiap hari.\n\n**2. Manajemen Stres:**\n\n* **Identifikasi pemicu stres:** Kenali hal-hal yang memicu stres Anda dan cari cara untuk menghindarinya atau mengatasinya.\n* **Teknik relaksasi:**  Coba teknik relaksasi seperti meditasi, yoga, pernapasan dalam, atau mendengarkan musik yang menenangkan.\n* **Manajemen waktu:**  Atur waktu Anda dengan efektif untuk menghindari perasaan terbebani.  Buat daftar prioritas dan fokus pada satu tugas pada satu waktu.\n* **Batasi paparan berita negatif:** Terlalu banyak terpapar berita negatif dapat meningkatkan kecemasan dan stres.\n\n**3. Hubungan Sosial yang Positif:**\n\n* **Jalin hubungan yang sehat:** Habiskan waktu bersama orang-orang yang mendukung dan peduli pada Anda.\n* **Berkomunikasi secara terbuka:**  Bicarakan perasaan Anda dengan orang yang Anda percaya.\n* **Bergabung dengan komunitas:** Temukan komunitas yang memiliki minat yang sama dengan Anda.\n* **Jangan ragu untuk meminta bantuan:**  Jika Anda merasa kesulitan, jangan ragu untuk meminta bantuan dari teman, keluarga, atau profesional.\n\n**4.  Aktivitas yang Menyenangkan:**\n\n* **Luangkan waktu untuk hobi:**  Lakukan aktivitas yang Anda sukai dan membuat Anda merasa bahagia.\n* **Pelajari hal baru:**  Mengasah keterampilan baru dapat meningkatkan rasa percaya diri dan memberikan kepuasan.\n* **Habiskan waktu di alam:** Berada di alam terbuka dapat memberikan efek menenangkan dan mengurangi stres.\n\n**5.  Perawatan Diri:**\n\n* **Prioritaskan kebutuhan Anda:**  Luangkan waktu untuk diri sendiri dan lakukan hal-hal yang membuat Anda merasa baik.\n* **Berpikir positif:**  Fokus pada hal-hal positif dalam hidup Anda dan latih rasa syukur.\n* **Tetapkan batasan yang sehat:**  Jangan ragu untuk mengatakan \"tidak\" pada hal-hal yang membuat Anda merasa tidak nyaman atau terbebani.\n\n**6.  Mencari Bantuan Profesional:**\n\n* **Jangan ragu untuk mencari bantuan:**  Jika Anda merasa kesulitan mengatasi masalah kesehatan mental Anda sendiri, jangan ragu untuk mencari bantuan profesional.\n* **Terapis atau psikolog:**  Mereka dapat membantu Anda mengidentifikasi masalah, mengembangkan strategi koping, dan memberikan dukungan emosional.\n* **Psikiater:**  Mereka dapat mendiagnosis dan mengobati gangguan mental dengan obat-obatan.\n\n\nIngatlah bahwa menjaga kesehatan mental adalah proses yang berkelanjutan.  Cobalah beberapa tips di atas dan temukan apa yang paling cocok untuk Anda.  Jika Anda merasa kesulitan, jangan ragu untuk mencari bantuan profesional.\n', '2025-04-18 17:28:15'),
(13, 14, 'Menjaga kesehatan mental sama pentingnya dengan menjaga kesehatan fisik. Berikut beberapa cara untuk menjaga kesehatan mental:\n\n**1. Pola Hidup Sehat:**\n\n* **Tidur yang Cukup:** Usahakan tidur 7-9 jam setiap malam. Kualitas tidur yang baik sangat penting untuk fungsi otak dan regulasi emosi.\n* **Pola Makan Seimbang:** Konsumsi makanan bergizi seimbang, batasi asupan gula, kafein, dan alkohol.\n* **Olahraga Teratur:** Aktivitas fisik melepaskan endorfin yang dapat meningkatkan suasana hati dan mengurangi stres.  Cukup 30 menit olahraga ringan hingga sedang beberapa kali seminggu.\n* **Paparan Sinar Matahari:** Sinar matahari membantu tubuh memproduksi vitamin D yang berperan penting dalam kesehatan mental.  Berjemurlah di pagi hari selama 10-15 menit.\n\n**2. Manajemen Stres:**\n\n* **Identifikasi Pemicu Stres:** Kenali hal-hal yang memicu stres Anda dan cari cara untuk menghindarinya atau mengelolanya.\n* **Teknik Relaksasi:** Latih teknik relaksasi seperti meditasi, yoga, pernapasan dalam, atau mindfulness.\n* **Manajemen Waktu:**  Atur waktu Anda dengan efektif untuk menghindari  kelebihan beban dan tekanan. Prioritaskan tugas dan istirahatlah secara teratur.\n\n**3. Hubungan Sosial yang Sehat:**\n\n* **Jalin Koneksi:**  Luangkan waktu untuk bersama orang-orang yang Anda cintai dan perkuat hubungan sosial Anda.\n* **Berbagi Perasaan:**  Jangan ragu untuk berbicara dengan teman, keluarga, atau terapis tentang apa yang Anda rasakan.\n* **Bergabung dengan Komunitas:**  Bergabung dengan komunitas yang memiliki minat yang sama dapat memberikan rasa kebersamaan dan dukungan.\n\n**4.  Kesehatan Mental Aktif:**\n\n* **Latih Rasa Syukur:**  Fokus pada hal-hal positif dalam hidup Anda dan luangkan waktu untuk mensyukurinya.\n* **Berpikir Positif:**  Ubah pola pikir negatif menjadi positif.\n* **Tetapkan Batasan:**  Belajar untuk mengatakan \"tidak\" terhadap hal-hal yang membuat Anda merasa tidak nyaman atau kewalahan.\n* **Kembangkan Hobi:**  Lakukan aktivitas yang Anda nikmati untuk mengisi waktu luang dan meningkatkan suasana hati.\n* **Batasi Paparan Berita Negatif:**  Terlalu banyak terpapar berita negatif dapat memengaruhi kesehatan mental Anda.\n\n**5.  Cari Bantuan Profesional:**\n\n* **Jangan Ragu untuk Mencari Bantuan:**  Jika Anda merasa kesulitan mengatasi masalah kesehatan mental Anda sendiri, jangan ragu untuk mencari bantuan profesional dari psikolog atau psikiater.  Ini bukan tanda kelemahan, melainkan langkah yang bijaksana untuk menjaga kesehatan mental Anda.\n\n**Ingat**:  Menjaga kesehatan mental adalah proses yang berkelanjutan.  Teruslah berlatih dan temukan strategi yang paling efektif untuk Anda.  Jangan ragu untuk mencoba hal-hal baru dan sesuaikan pendekatan Anda seiring waktu.\n', '2025-04-18 17:31:29');

-- --------------------------------------------------------

--
-- Table structure for table `sesi_konsultasi`
--

CREATE TABLE `sesi_konsultasi` (
  `sesi_id` int(11) NOT NULL,
  `user_id` int(11) NOT NULL,
  `status` enum('Berlangsung','Selesai','Dibatalkan') NOT NULL DEFAULT 'Berlangsung',
  `waktu_mulai` datetime DEFAULT current_timestamp(),
  `waktu_selesai` datetime DEFAULT NULL,
  `auto_forward` tinyint(1) DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `sesi_konsultasi`
--

INSERT INTO `sesi_konsultasi` (`sesi_id`, `user_id`, `status`, `waktu_mulai`, `waktu_selesai`, `auto_forward`) VALUES
(4, 7, 'Berlangsung', '2025-03-30 08:38:00', NULL, 0),
(5, 9, 'Berlangsung', '2025-04-18 13:45:47', NULL, 0),
(6, 9, 'Berlangsung', '2025-04-18 13:46:12', NULL, 0),
(7, 9, 'Berlangsung', '2025-04-18 13:54:48', NULL, 0);

-- --------------------------------------------------------

--
-- Table structure for table `user`
--

CREATE TABLE `user` (
  `user_id` int(11) NOT NULL,
  `nama` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `password` varchar(255) NOT NULL,
  `role` enum('admin','user') NOT NULL DEFAULT 'user'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;

--
-- Dumping data for table `user`
--

INSERT INTO `user` (`user_id`, `nama`, `email`, `password`, `role`) VALUES
(7, 'Ilham Updated', 'ilham.new@example.com', 'newpassword123', 'user'),
(9, 'admin', 'admin@gmail.com', 'admin', 'admin');

--
-- Indexes for dumped tables
--

--
-- Indexes for table `konsultan_kontak`
--
ALTER TABLE `konsultan_kontak`
  ADD PRIMARY KEY (`konsultan_id`);

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
  ADD KEY `user_id` (`user_id`);

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
-- AUTO_INCREMENT for table `konsultan_kontak`
--
ALTER TABLE `konsultan_kontak`
  MODIFY `konsultan_id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=2;

--
-- AUTO_INCREMENT for table `pesan_konsultasi`
--
ALTER TABLE `pesan_konsultasi`
  MODIFY `pesan_id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=15;

--
-- AUTO_INCREMENT for table `respons_chatbot`
--
ALTER TABLE `respons_chatbot`
  MODIFY `respons_id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=14;

--
-- AUTO_INCREMENT for table `sesi_konsultasi`
--
ALTER TABLE `sesi_konsultasi`
  MODIFY `sesi_id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=8;

--
-- AUTO_INCREMENT for table `user`
--
ALTER TABLE `user`
  MODIFY `user_id` int(11) NOT NULL AUTO_INCREMENT, AUTO_INCREMENT=11;

--
-- Constraints for dumped tables
--

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
  ADD CONSTRAINT `sesi_konsultasi_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `user` (`user_id`) ON DELETE CASCADE;
COMMIT;

/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
