# TFL Bus Bunching Detector (T-UGS)

## Description

T-UGS is a command-line tool written in Go that analyzes live Transport for London (TfL) bus arrival data to detect "bus bunching." Bus bunching occurs when two or more buses on the same route arrive at a stop in very close succession, disrupting the expected service headway.

This tool can fetch live arrival data for specific bus stops, analyze the time between arrivals for each line, and report any instances of bunching based on a configurable time threshold.

## Features

* Fetches real-time bus arrival predictions from the official TfL Unified API.

* Parses and loads a list of all 20,000+ London bus stops from a local CSV file for searching.

* Allows users to search for bus stops by name in a case-insensitive way.

* Provides an interactive menu to select one or more stops from the search results.

* Analyzes arrival data to identify buses on the same line that are arriving closer than a specified time threshold.

* Displays a clean, user-friendly summary of any bunching events found.

## Prerequisites

* Go (version 1.19 or later recommended)

* A valid TfL API Key. You can get one from the [TfL API Portal](https://api-portal.tfl.gov.uk/).

## Installation & Setup

1. **Clone the repository:**

   ```bash
   git clone <your-repository-url>
   cd tfl-bunching-detector
   ```

2. **Set up your API Key:**
   The application loads your TfL API key from a `.env` file at the root of the project. Create this file:

   ```bash
   touch .env
   ```

   Then, open the `.env` file and add your API key on a single line:

   ```
   TFL_API_KEY=your_tfl_api_key_goes_here
   ```

   **Important:** Make sure your `.gitignore` file contains `.env` to prevent accidentally committing your API key.

3. **Tidy dependencies:**
   Run the following command to ensure all package dependencies are correct:

   ```bash
   go mod tidy
   ```

## Usage

The application can be run in two main modes: searching for a stop or analyzing a known stop ID directly.

### 1. Search Mode

This is the most common way to use the tool. You provide a search term, and the application will present a list of matching stops for you to choose from.

**Command:**

```bash
go run ./cmd/tugs -search="<search_term>"
```

**Example:**

```bash
go run ./cmd/tugs -search="Victoria Station"
```

This will output a numbered list of all stops containing "Victoria Station" in their name. You will then be prompted to enter the number of the stop (or stops) you wish to analyze. The tool supports selecting multiple stops by entering numbers separated by a space (e.g., `1 3 5`).

### 2. Direct ID Mode

If you already know the Naptan ID of the bus stop you want to analyze, you can provide it directly.

**Command:**

```bash
go run ./cmd/tugs -stop-id="<naptan_id>"
```

**Example:**

```bash
go run ./cmd/tugs -stop-id="490000234H"
```

### Optional Flags

You can combine the following optional flag with either of the modes above.

* **`-threshold`**: Sets the time in seconds to consider as bunching. If the time between two buses is less than this value, it will be reported. Defaults to `90`.

**Example with custom threshold:**

```bash
go run ./cmd/tugs -search="Waterloo" -threshold=120