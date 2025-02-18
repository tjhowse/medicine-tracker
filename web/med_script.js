const medicineLog = document.getElementById('medicine-log-table');
const medicineLogTable = document.createElement('table');
const medicineLogHeaderRow = document.createElement('tr');
medicineLogHeaderRow.innerHTML = '<th>Time</th><th>Count</th><th>Medicine</th><th>Note</th>';
populateMedicineLogTab();

const availableMedicines = document.getElementById('available-medicines');
// Populate a table into this div
const medicinesTable = document.createElement('table');
const medicinesHeaderRow = document.createElement('tr');
medicinesHeaderRow.innerHTML = '<th>ID</th><th>Type</th><th>Dose (mg)</th>';
medicinesTable.appendChild(medicinesHeaderRow);
updateMedicinesTable();
availableMedicines.appendChild(medicinesTable);

populateSettings();

openDefaultTab();

function populateMedicineLogTab() {
  // Retrieve today's workout details and populate the Today's Workout section

  fetch('/api/v1/api/v1/medicine-log', {
    headers: {
      'Content-Type': 'application/json',
    },
  })
  .then(response => response.json())
  .then(data => {
    // Clear the table
    while (medicineLogTable.rows.length > 1) {
      medicineLogTable.deleteRow(1);
    }
    data.forEach(detail => {
      addMedicineLogRow(medicineLogTable, detail.time, detail.count, detail.medicine_id, detail.note);
    });

  }).catch(error => {
    console.error(error);
    window.location.href = "/login.html";
  });
}

function addMedicineLogRow(table, time, count, medicine_id, note) {
  const row = document.createElement('tr');
  const timeCell = document.createElement('td');
  const countCell = document.createElement('td');
  const medicineCell = document.createElement('td');
  const noteCell = document.createElement('td');

  timeCell.textContent = time;
  countCell.textContent = count;
  medicineCell.textContent = medicine_id;
  noteCell.textContent = note;

  row.appendChild(timeCell);
  row.appendChild(countCell);
  row.appendChild(medicineCell);
  row.appendChild(noteCell);
  table.appendChild(row);
}

function updateMedicinesTable() {

  fetch('/api/v1/api/v1/medicines', {
    headers: {
      'Content-Type': 'application/json',
    },
  })
  .then(response => response.json())
  .then(data => {
    populateMedicinesTable(data);
  }).catch(error => {
    console.error(error);
  });
}

function populateMedicinesTable(data) {
  // Clear the table
  while (medicinesTable.rows.length > 1) {
    medicinesTable.deleteRow(1);
  }

  max_id = 0;
  data.forEach(detail => {
    addMedicinesRow(medicinesTable, detail.medicine_id, detail.name, detail.dose);
    if (detail.medicine_id > max_id) {
      max_id = detail.medicine_id;
    }
  });
  addMedicinesRow(medicinesTable, max_id + 1, "", "");
  medicinesTable.rows[medicinesTable.rows.length-1].cells[1].getElementsByTagName('input')[0].focus();

}

function saveMedicines() {
  const medicines = [];
  for (i = 1; i < medicinesTable.rows.length; i++) {
    const row = medicinesTable.rows[i];
    const id = parseFloat(row.cells[0].textContent);
    const name = row.cells[1].getElementsByTagName('input')[0].value;
    const dose = parseFloat(row.cells[2].getElementsByTagName('input')[0].value);
    if (name === "" || isNaN(dose)) {
      continue;
    }
    medicines.push({medicine_id: id, name: name, dose: dose});
  }

  console.log(medicines);
  fetch('/api/v1/api/v1/medicines', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(medicines),
  })
  .then(response => response.json())
  .then(data => {
    populateMedicinesTable(data);
  }).catch(error => {
    console.error(error);
  });
}

function logout() {
  fetch('/api/v1/api/v1/logout')
  .then(() => {
    window.location.href = "/login.html";
  });
}

function deleteUser() {
  // Show an alert prompt to confirm
  if (!confirm("Are you sure you want to delete your account? There is absolutely no way to recover it.")) {
    return;
  }
  console.log("Deleting user");
  fetch('/api/v1/api/v1/delete-user')
  .then(() => {
    window.location.href = "/login.html";
  });
}

function populateSettings() {
  fetch('/api/v1/api/v1/settings', {
    headers: {
      'Content-Type': 'application/json',
    },
  })
  .then(response => response.json())
  .then(data => {
    if (!data) {
      return;
    }

    // Populate the Greeting div
    var greetingDiv = document.getElementById('Greeting');
    greetingDiv.textContent = "Hello, " + data.name + "!";
  }).catch(error => {
    console.error(error);
  });
}

function addMedicinesRow(table, id, name, dose) {
  const row = document.createElement('tr');
  const idCell = document.createElement('td');
  const nameCell = document.createElement('td');
  const doseCell = document.createElement('td');
  idCell.textContent = id;

  const nameInput = document.createElement('input');
  nameInput.onsubmit = saveMedicines;
  nameInput.value = name;
  nameCell.appendChild(nameInput);

  const doseInput = document.createElement('input');
  nameInput.onsubmit = saveMedicines;
  doseInput.value = dose;
  doseCell.appendChild(doseInput);

  row.appendChild(idCell);
  row.appendChild(nameCell);
  row.appendChild(doseCell);
  table.appendChild(row);
}



function openDefaultTab() {
  const cookie = document.cookie;
  const cookieParts = cookie.split(";");
  for (i = 0; i < cookieParts.length; i++) {
    const cookiePart = cookieParts[i];
    const cookieName = cookiePart.split("=")[0];
    if (cookieName.trim() === "tab") {
      tabName = cookiePart.split("=")[1];
      // Click the tab header
      document.getElementById("tab"+tabName).click();
      return;
    }
  }
  // If we didn't find a cookie, open the workout tab
  document.getElementById("tabMedicineLog").click();
}

function openTab(event, tabName) {
  // Declare all variables
  let i, tabcontent, tablinks;

  // Get all elements with class="tabcontent" and hide them
  tabcontent = document.getElementsByClassName("tabcontent");
  for (i = 0; i < tabcontent.length; i++) {
    tabcontent[i].style.display = "none";
  }

  // Get all elements with class="tablinks" and remove the class "active"
  tablinks = document.getElementsByClassName("tablinks");
  for (i = 0; i < tablinks.length; i++) {
    tablinks[i].className = tablinks[i].className.replace(" active", "");
  }

  // Show the current tab, and add an "active" class to the button that opened the tab
  document.getElementById(tabName).style.display = "block";
  event.currentTarget.className += " active";

  // Set a cookie to remember which tab we had selected on load
  document.cookie = "tab=" + tabName;

}

function exportMedicineLogDataToCSV() {
  fetch('/api/v1/api/v1/medicine-log?all=true', {
    headers: {
      'Content-Type': 'application/json',
    },
  })
  .then(response => response.json())
  .then(data => {
    let csvContent = "data:text/csv;charset=utf-8," + "time,count,dose,type,note\n"
    data.forEach(detail => {
      csvContent += detail.time+","+detail.count + "," + detail.medicine_id +",,\n";
      console.log(csvContent);
    });
    var encodedUri = encodeURI(csvContent);
    window.open(encodedUri);
  }).catch(error => {
    console.error(error);
  });
}
