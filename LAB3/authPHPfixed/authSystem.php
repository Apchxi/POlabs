<?php

// Prepared statements prevent SQL injection vulnerabilities
if (isset($_GET['Login'])) {
  // Sanitize user input to prevent XSS vulnerabilities
  $username = mysqli_real_escape_string($GLOBALS["___mysqli_ston"], $_GET['username']);
  // Password hashing - use a secure hashing algorithm like bcrypt
  $password = password_hash($_GET['password'], PASSWORD_BCRYPT);

  // Prepare the query
  $query = $GLOBALS["___mysqli_ston"]->prepare("SELECT * FROM `users` WHERE user = ? AND password = ?");
  // Bind the parameters
  $query->bind_param("ss", $username, $password);
  // Execute the query
  $query->execute();
  // Fetch the results
  $result = $query->get_result();

  // Check login success
  if ($result->num_rows === 1) {
    // Get user details
    $row = $result->fetch_assoc();
    $avatar = $row["avatar"];

    // Login successful
    $html .= "<p>Welcome to the password protected area {$username}</p>";
    $html .= "<img src=\"{$avatar}\" />";
  } else {
    // Login failed
    $html .= "<pre><br />Username and/or password incorrect.</pre>";
  }

  // Close the prepared statement and connection
  $query->close();
  mysqli_close($GLOBALS["___mysqli_ston"]);
}
?>