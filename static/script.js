document.getElementById("reportForm").addEventListener("submit", function (e) {
  e.preventDefault();

  // Simulate API call. We will change this after building the api from which we will get the information from
  setTimeout(() => {
    const success = Math.random() > 0.1; // 90% success rate simulation

    if (success) {
      document.getElementById("successMessage").style.display = "block";
      document.getElementById("errorMessage").style.display = "none";
      document.getElementById("reportForm").reset();

      setTimeout(() => {
        document.getElementById("successMessage").style.display = "none";
      }, 5000);
    } else {
      document.getElementById("errorMessage").style.display = "block";
      document.getElementById("successMessage").style.display = "none";

      setTimeout(() => {
        document.getElementById("errorMessage").style.display = "none";
      }, 5000);
    }
  }, 1500);
});

// Set minimum date-time to 30 days ago
const dateTimeInput = document.getElementById("incidentDateTime");
const thirtyDaysAgo = new Date();
thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30);
dateTimeInput.max = new Date().toISOString().slice(0, 16);
dateTimeInput.min = thirtyDaysAgo.toISOString().slice(0, 16);
