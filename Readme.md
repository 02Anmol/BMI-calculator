# BMI Calculator Web Application

A simple web-based BMI (Body Mass Index) calculator built with Go that allows users to calculate and track BMI records.

## Features

- Calculate BMI based on weight (kg) and height (m)
- Categorize BMI results (Underweight, Normal Weight, Overweight, Obesity)
- Store user records persistently in a JSON file
- View all calculated BMI records in a table
- Clean web interface using HTML templates

## Prerequisites

- Go 1.16 or higher installed on your system
- Basic understanding of running Go applications

## Project Structure
```
.
├── main.go                 # Main application file
├── users_data.json         # Data storage file (auto-created)
└── templates/
    └── *.html             # HTML template files
```

## Installation

1. Clone or download this repository
2. Navigate to the project directory:
```bash
   cd bmi-calculator
```

3. Ensure you have a `templates` folder with your HTML template files

## Configuration

The application uses a JSON file for data storage. You need to define the data file path in your code:
```go
const dataFile = "users_data.json"
```

**Note:** Uncomment this line in `main.go` before running the application.

## Running the Application

1. Start the server:
```bash
   go run main.go
```

2. Open your web browser and navigate to:
```
   http://localhost:8080
```

3. The application will start on port 8080 and display logs in the terminal

## Usage

1. **Enter User Information:**
   - Name: Enter the person's name
   - Weight: Enter weight in kilograms (positive numbers only)
   - Height: Enter height in meters (positive numbers only)

2. **Calculate BMI:**
   - Click the submit button
   - The BMI will be calculated automatically
   - Results are categorized into health ranges

3. **View Records:**
   - All calculated BMI records are displayed in a table
   - Records persist between server restarts

## BMI Categories

The application categorizes BMI results as follows:

| BMI Range | Category |
|-----------|----------|
| < 18.5 | Underweight |
| 18.5 - 24.9 | Normal Weight |
| 25.0 - 29.9 | Overweight |
| ≥ 30.0 | Obesity |

## Data Storage

- User records are stored in `users_data.json`
- Data persists between application restarts
- If the file doesn't exist on first run, it will be created automatically
- Each record contains: name, weight, height, calculated BMI, and category

## API Endpoints

- `GET /` - Display the main page with form and records table
- `POST /calculate` - Process form submission and calculate BMI

## Error Handling

The application handles:
- Missing data file (starts with empty list)
- Invalid input validation (negative or non-numeric values)
- File read/write errors
- Template rendering errors

## Troubleshooting

**Server won't start:**
- Check if port 8080 is already in use
- Ensure Go is properly installed (`go version`)

**Templates not found:**
- Verify the `templates/` folder exists
- Check that HTML files are present in the templates folder

**Data not persisting:**
- Ensure the `dataFile` constant is uncommented
- Check file write permissions in the directory

**Invalid input errors:**
- Weight and height must be positive numbers
- Use decimal point for fractional values (e.g., 1.75 for height)

## Example Data Format

The `users_data.json` file stores data in this format:
```json
[
  {
    "name": "John Doe",
    "weight_kg": 75.5,
    "height_m": 1.75,
    "bmi": 24.65,
    "category": "Normal Weight"
  }
]
```

## Development

To modify the application:

1. **Change the port:** Edit the `port` variable in `main()`
2. **Modify BMI categories:** Update the `getBMICategory()` function
3. **Customize templates:** Edit HTML files in the `templates/` folder
4. **Change data storage:** Modify `loadUserData()` and `saveUserData()` functions

## License

This project is provided as-is for educational purposes.

## Contributing

Feel free to fork this project and submit pull requests for any improvements.

## Support

For issues or questions, please open an issue in the repository.