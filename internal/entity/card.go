package entity

type Card struct {
	ID              int64
	PatientID       int64
	ChronicDiseases string
	DisabilityGroup int
	BloodType       int
	RhFactor        bool
	Consultations   []Consultation
}

type Consultation struct {
	ID              int64
	CardID          int64
	DoctorID        string
	FullName        string
	Complaints      string
	Descriptions    string
	Recommendations string
}
