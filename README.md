# 🚌 TfL Bus Bunching Detector

This command-line tool analyses live London bus arrival data to detect and report instances of bus bunching on specific routes or at specific stops. It uses the Transport for London (TfL) API to fetch real-time data and applies a configurable threshold to identify when buses are arriving too close together (bunching).

## ✨ Features
- Detects bus bunching events for a given route or stop using live TfL data
- Allows searching for stops by name or Naptan ID
- Supports analysis of specific bus lines/routes
- Configurable bunching threshold (in seconds)
- Tabular output of bunching events, including line, stop, vehicle IDs, and headway

## 🚀 Usage

```
go run ./cmd/tugs/main.go [flags]
```

### 🏷️ Flags
- `-stop-id` (string): NaptanId for a specific stop (e.g. `-stop-id="490008660N"`)
- `-search` (string): Search for a stop by name (e.g. `-search="Waterloo"`)
- `-line` (string): Line ID or Line Number to analyse a whole route (e.g. `-line="77"`)
- `-threshold` (int): Threshold for bunched buses in seconds (default: 90)

At least one of `-stop-id`, `-search`, or `-line` must be provided.

### 💡 Examples
- Analyse bunching for a specific stop:
  ```
  go run ./cmd/tugs/main.go -stop-id="490008660N"
  ```
- Search for stops by name and select interactively:
  ```
  go run ./cmd/tugs/main.go -search="Waterloo"
  ```
- Analyse bunching for a whole bus route with a custom threshold:
  ```
  go run ./cmd/tugs/main.go -line="77" -threshold=500
  ```

## 📊 Output
- If bunching events are detected, a table is printed with columns: Line, Stop, Vehicle 1, Vehicle 2, Headway (s)
- If no bunching is detected, a message is printed

## 🛠️ Requirements
- Go 1.18+
- A valid TfL API key (set in your config)

## ⚙️ Configuration

### 🔑 Getting a TfL API Key
1. Visit the [TfL API Portal](https://api-portal.tfl.gov.uk/signup) and sign up for a free account.
2. After registering and verifying your email, log in and create a new application to obtain your API key (app_key).

### 📝 Setting the API Key
1. Create a file named `.env` in the project root directory (next to `main.go`).
2. Add the following line to your `.env` file, replacing `YOUR_TFL_API_KEY` with your actual key:
   ```
   TFL_APP_KEY=YOUR_TFL_API_KEY
   ```
3. The application will automatically load this key from the `.env` file at startup.

If you prefer, you can also set the environment variable directly in your shell before running the app:
```
export TFL_APP_KEY=YOUR_TFL_API_KEY
```

The app requires this key to authenticate requests to the TfL API.

## 🗂️ Bus Stop Data

The application requires a list of all London bus stops in CSV format.

### Downloading the Data
1. Download the latest bus stop data from the official TfL source:
   [https://tfl.gov.uk/bus-stops.csv](https://tfl.gov.uk/bus-stops.csv)
2. Save the downloaded file as:
   `internal/data/bus-stops.csv`
   (relative to the project root directory)

This file is used by the application to look up and search for bus stops by name or ID.

## 📄 License
MIT
