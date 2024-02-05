# Financial Tracker App

## Overview

This is a simple financial application written in Golang, designed to help users track their expenses and income conveniently. The application uses a SQLite database to store and manage financial data and is enhanced with HTMX for a seamless and dynamic user experience.

## Features

- **Expense Tracking:** Easily log your daily expenses and income transactions.
- **SQLite Database:** Utilizes a lightweight SQLite database for data storage.
- **HTMX Integration:** Enhances the user interface with dynamic content updates, providing a smooth and responsive experience.

## Getting Started

### Prerequisites

Make sure you have the following installed:

- [Golang](https://golang.org/dl/): The programming language used for the application.
- [SQLite](https://www.sqlite.org/download.html): A C library that provides a lightweight disk-based database.
- [HTMX](https://htmx.org/): A library for modern web applications.

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/your-username/financial-tracker.git
   ```

2. Change directory to the project folder:

   ```bash
   cd financial-tracker
   ```

3. Set up the SQLite database:

   ```bash
   go run setup.go
   ```

4. Build and run the application:

   ```bash
   go build && ./financial-tracker
   ```

5. Open your web browser and navigate to [http://localhost:8080](http://localhost:8080) to access the application.

## Usage

- Add Expenses: Log your expenses and income transactions.
- View Transactions: Access a summary of your financial transactions.
- Analyze Trends: Use the application to analyze your spending patterns over time.

## Contributing

If you would like to contribute to this project, feel free to fork the repository and submit a pull request. Please follow the [Contribution Guidelines](CONTRIBUTING.md).

## License

This project is licensed under the [MIT License](LICENSE).

## Acknowledgments

- Thanks to the Golang community for providing a powerful and efficient programming language.
- HTMX for enhancing the web application's user interface.

## Contact

For any issues or inquiries, please contact [your-email@example.com](mailto:your-email@example.com).

Happy tracking!
