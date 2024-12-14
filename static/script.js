// Set minimum date-time to 30 days ago
const dateTimeInput = document.getElementById("incidentDateTime");
const thirtyDaysAgo = new Date();
thirtyDaysAgo.setDate(thirtyDaysAgo.getDate() - 30);
dateTimeInput.max = new Date().toISOString().slice(0, 16);
dateTimeInput.min = thirtyDaysAgo.toISOString().slice(0, 16);
