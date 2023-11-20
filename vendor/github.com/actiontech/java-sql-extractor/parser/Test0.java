import java.sql.Connection;
import java.sql.ResultSet;
import java.sql.SQLException;
import java.sql.Statement;

public class DatabaseConnection {
    private Connection connection;


    public Connection getConnection() {
        return connection;
    }

    public ResultSet executeQuery(String sqlQuery, String sqlQuery1) throws SQLException {
        Statement statement = connection.createStatement();
        String a = sqlQuery+sqlQuery1;
        ResultSet resultSet = statement.executeQuery(a);
        return resultSet;
    }
}

public class DatabaseQuery {
    private Connection connection;

    public DatabaseQuery(Connection connection) {
        this.connection = connection;
    }



    public static void main(String[] args) {
        String jdbcUrl = "jdbc:mysql://localhost:3306/mydatabase";
        String username = "yourUsername";
        String password = "yourPassword";

        try {
            DatabaseConnection dbQuery = new DatabaseConnection();

            String sqlQuery = "SELECT * FROM employees";
            String sqlQuery1 = " limit 1";
            ResultSet resultSet = dbQuery.executeQuery(sqlQuery, sqlQuery1);



            connection.close();
        } catch (SQLException e) {
            e.printStackTrace();
        }
    }
}

public class DatabaseQuery1 {
    private Connection connection;

    public DatabaseQuery1(Connection connection) {
        this.connection = connection;
    }



    public static void main(String[] args) {
        String jdbcUrl = "jdbc:mysql://localhost:3306/mydatabase";
        String username = "yourUsername";
        String password = "yourPassword";

        try {
            DatabaseConnection dbQuery = new DatabaseConnection();

            String sqlQuery = "SELECT * FROM users";
            String sqlQuery1 = " limit 1";
            ResultSet resultSet = dbQuery.executeQuery(sqlQuery, sqlQuery1);



            connection.close();
        } catch (SQLException e) {
            e.printStackTrace();
        }
    }
}
