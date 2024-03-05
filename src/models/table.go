package models

import "time"

type User struct {
	ID             int    `gorm:"primaryKey;autoIncrement"`
	Nama           string `gorm:"not null;size:255"`
	Email          string `gorm:"not null;unique;size:255"`
	Password       string `gorm:"not null;size:255"`
	Role           string `gorm:"not null;size:50"`
	JawabanPeserta []JawabanPeserta
}

type Quiz struct {
	ID             int64     `gorm:"primaryKey;autoIncrement"`
	Judul          string    `gorm:"not null;size:255"`
	Deskripsi      string    `gorm:"type:text"`
	WaktuMulai     time.Time `gorm:"not null"`
	WaktuSelesai   time.Time `gorm:"not null"`
	Pertanyaan     []Pertanyaan
	JawabanPeserta []JawabanPeserta
}

type Pertanyaan struct {
	ID             int    `gorm:"primaryKey;autoIncrement"`
	Pertanyaan     string `gorm:"type:text;not null"`
	OpsiJawaban    string `gorm:"type:text;not null"`
	JawabanBenar   int    `gorm:"not null"`
	QuizID         int    `gorm:"not null"`
	JawabanPeserta []JawabanPeserta
}

type JawabanPeserta struct {
	ID             int `gorm:"primaryKey;autoIncrement"`
	UserID         int `gorm:"not null"`
	QuizID         int `gorm:"not null"`
	PertanyaanID   int `gorm:"not null"`
	JawabanPeserta int `gorm:"not null"`
	Skor           int
}
