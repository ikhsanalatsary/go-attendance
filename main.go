package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

var version = "4.3.0"

// Auth object, it's for login request the structure like this
type Auth struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Data object, it's a response as part of Response object...
type Data struct {
	Token                    string `json:"token,omitempty"`
	CurrentShiftDate         string `json:"current_shift_date,omitempty"`
	CurrentShiftID           int64  `json:"current_shift_id,omitempty"`
	CurrentShiftName         string `json:"current_shift_name,omitempty"`
	CurrentShiftScIn         string `json:"current_shift_sc_in,omitempty"`
	CurrentShiftScOut        string `json:"current_shift_sc_out,omitempty"`
	CurrentShiftScBreakStart string `json:"current_shift_sc_break_start,omitempty"`
	CurrentShiftScBreakEnd   string `json:"current_shift_sc_break_end,omitempty"`
	IsCheckIn                bool   `json:"is_check_in,omitempty"`
	ActualCheckIn            string `json:"actual_check_in,omitempty"`
	IsCheckOut               bool   `json:"is_check_out,omitempty"`
	ActualCheckOut           string `json:"actual_check_out,omitempty"`
	IsBreakStart             bool   `json:"is_break_start,omitempty"`
	ActualBreakStart         string `json:"actual_break_start,omitempty"`
	IsBreakEnd               bool   `json:"is_break_end,omitempty"`
	ActualBreakEnd           string `json:"actual_break_end,omitempty"`
	ServerTime               string `json:"server_time,omitempty"`
	CountLocation            int64  `json:"count_location,omitempty"`
	AttendanceType           string `json:"attendanceType,omitempty"`
	IsLate                   string `json:"isLate,omitempty"`
}

// Response object. Either login or attendance use this
type Response struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
	Data    *Data  `json:"data"`
	Version string `json:"version"`
}

// Attendance object is a request for hit attendance. the structure like this
type Attendance struct {
	Latitude    string `json:"latitude"`
	Longitude   string `json:"longitude"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func main() {
	port := os.Getenv("PORT")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			fmt.Println("YES")
		}
		w.Write([]byte("Hello Go!"))
	})

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			dec := json.NewDecoder(r.Body)
			var body Auth
			err := dec.Decode(&body)
			w.Header().Add("Content-Type", "application/json")
			if err != nil {
				w.WriteHeader(400)
				return
			}
			resp := Response{}
			resp.Message = "Successfully login"
			resp.Status = 200
			resp.Version = version
			if body.Email == "" || body.Password == "" {
				resp.Message = "Email or Password missing"
				resp.Status = 400
			} else {
				data := &Data{}
				data.Token = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJyYWZhcmNAd2FsbGV4dGVjaC5jb20iLCJpYXQiOjE2MDkzNDE0NTksImV4cCI6MTYwOTM3MDI1OX0=.vsXfCp62cs8DyYU0-peVqOAolBd4Xmn1SUbL1cQxnjM="
				resp.Data = data
			}
			// fmt.Printf("%#v", body.Password)
			bts, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(400)
				return
			}
			w.WriteHeader(resp.Status)
			w.Write(bts)
			return
		}
		w.WriteHeader(404)
		return
	})

	http.HandleFunc("/live-attendance", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost {
			token := r.Header.Get("Authorization")
			fmt.Printf("%#v", token)
			resp := Response{}
			resp.Status = 200
			resp.Version = version
			dec := json.NewDecoder(r.Body)
			var body Attendance
			err := dec.Decode(&body)
			w.Header().Add("Content-Type", "application/json")
			if err != nil {
				w.WriteHeader(400)
				return
			}
			lat := body.Latitude
			long := body.Longitude
			desc := body.Description
			status := body.Status
			if token == "" {
				resp.Message = "Credential missing!"
				resp.Status = 401
			} else {
				resp.Message = "Successfully " + body.Status
				resp.Data = &Data{IsCheckIn: body.Status == "checkin", IsCheckOut: body.Status == "checkout", IsLate: "Not Late"}
			}
			fmt.Println(lat, long, desc, status)
			bts, err := json.Marshal(resp)
			if err != nil {
				w.WriteHeader(400)
				return
			}
			w.WriteHeader(resp.Status)
			w.Write(bts)
			return
		}
		w.WriteHeader(404)
		return
	})

	log.Fatal(http.ListenAndServe(port, nil))
}
