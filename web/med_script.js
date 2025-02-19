
const medicineLog = document.getElementById('medicine-log-table');
const medicineLogTable = document.createElement('table');
const medicineLogHeaderRow = document.createElement('tr');
medicineLogHeaderRow.innerHTML = '<th>Time</th><th>Count</th><th>Medicine</th><th>Note</th><th>Delete</th>';
medicineLogTable.appendChild(medicineLogHeaderRow);
medicineLog.appendChild(medicineLogTable);
populateMedicineLogTable();

const medicineLogForm = document.getElementById('medicine-log-form');
// Add a header to the form
const medicineLogFormHeader = document.createElement('h2');
medicineLogFormHeader.textContent = "Log Medicine";
medicineLogForm.appendChild(medicineLogFormHeader);
const datetimeInput = document.createElement('input');
datetimeInput.type = 'datetime-local';
const now = new Date();
const tzOffset = now.getTimezoneOffset() * 60000; // offset in milliseconds
const localISOTime = new Date(now - tzOffset).toISOString().slice(0, 16);
datetimeInput.value = localISOTime;
medicineLogForm.appendChild(datetimeInput);
// Add a row of buttons with 0.5, 1 an 2 on them that will set the count input
const countButtons = document.createElement('div');
countButtons.id = 'count-buttons';
const counts = [0.5, 1, 2];
counts.forEach(count => {
  const button = document.createElement('button');
  button.textContent = count;
  button.onclick = function() {
    countInput.value = count;
  };
  countButtons.appendChild(button);
});
medicineLogForm.appendChild(countButtons);
// Add a count input to the form
const countInput = document.createElement('input');
countInput.type = 'float';
countInput.placeholder = 'Count';
// Give the countInput the focus when the form is shown
medicineLogForm.appendChild(countInput);
// Add a medicine select to the form
const medicineSelect = document.createElement('select');
medicineSelect.id = 'medicine-select';
medicineLogForm.appendChild(medicineSelect);
populateMedicineDropdown(medicineSelect);

// Add a note input to the form
const noteInput = document.createElement('input');
noteInput.type = 'text';
noteInput.placeholder = 'Note';
medicineLogForm.appendChild(noteInput);
// Add a submit button to the form
submitButton = document.createElement('button');
submitButton.textContent = 'Submit';
submitButton.onclick = function() {
  // Just hardcode the timezone offset:
  // isoTZSuffix = "+10:00";

  // const time = datetimeInput.value + isoTZSuffix;
  const time = datetimeInput.value+":00+10:00";
  const count = parseFloat(countInput.value);
  const medicine_id = parseFloat(medicineSelect.value);
  const note = noteInput.value;

  fetch('/api/v1/api/v1/medicine-log', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({time: time, count: count, medicine_id: medicine_id, note: note}),
  })
  .then(response => response.json())
  .then(data => {
    populateMedicineLogTable();
    countInput.value = "";
    medicineSelect.value = -1;
    noteInput.value = "";
    // Put focus on the value field
    countInput.focus();
  }).catch(error => {
    console.error(error);
    window.location.href = "/login.html";
  });
};
medicineLogForm.appendChild(submitButton);
medicineLogForm.onsubmit = function(event) {
  event.preventDefault();
};

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

function populateMedicineDropdown(select) {
  const option = document.createElement('option');
  option.value = -1;
  option.textContent = "N/A";
  select.appendChild(option);
  fetch('/api/v1/api/v1/medicines', {
    headers: {
      'Content-Type': 'application/json',
    },
  })
  .then(response => response.json())
  .then(data => {
    data.forEach(detail => {
      const option = document.createElement('option');
      option.value = detail.medicine_id;
      option.textContent = detail.name+" "+detail.dose +"mg";
      select.appendChild(option);
    });
  }).catch(error => {
    console.error(error);
    window.location.href = "/login.html";
  });
}


function populateMedicineLogTable() {
  // Retrieve the medicine log and populate the table


  fetch('/api/v1/api/v1/medicine-log', {
    headers: {
      'Content-Type': 'application/json',
    },
  })
  .then(response => response.json())
  .then(logData => {
    // Clear the table
    while (medicineLogTable.rows.length > 1) {
      medicineLogTable.deleteRow(1);
    }
    // get the name map
    fetch('/api/v1/api/v1/medicines', {
      headers: {
        'Content-Type': 'application/json',
      },
    })
    .then(response => response.json())
    .then(data => {
      const medicineIDtoNameMap = {};
      data.forEach(detail => {
        medicineIDtoNameMap[detail.medicine_id] = detail.name;
      });

      logData.forEach(detail => {
        console.log(detail);
        addMedicineLogRow(medicineLogTable, detail.log_id, detail.time, detail.count, medicineIDtoNameMap[detail.medicine_id], detail.note);
      });
    }).catch(error => {
      console.error(error);
    });

  }).catch(error => {
    console.error(error);
    window.location.href = "/login.html";
  });
}

function addMedicineLogRow(table, id, time, count, medicine_id, note) {
  const row = document.createElement('tr');
  const timeCell = document.createElement('td');
  const countCell = document.createElement('td');
  const medicineCell = document.createElement('td');
  const noteCell = document.createElement('td');
  const deleteButtonCell = document.createElement('td');

  const formattedTime = time.slice(0, 16).replace('T', ' ');
  timeCell.textContent = formattedTime;
  countCell.textContent = count;
  medicineCell.textContent = medicine_id;
  noteCell.textContent = note;
  deleteButton = document.createElement('button');
  deleteButton.textContent = "Delete";
  deleteButton.onclick = function() {
    fetch(`/api/v1/api/v1/medicine-log?log_id=${id}`, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
      },
        })
        .then(response => response.json())
        .then(data => {
      populateMedicineLogTable();
        }).catch(error => {
      console.error(error);
      window.location.href = "/login.html";
        });
  };
  deleteButtonCell.appendChild(deleteButton);

  row.appendChild(timeCell);
  row.appendChild(countCell);
  row.appendChild(medicineCell);
  row.appendChild(noteCell);
  row.appendChild(deleteButtonCell);
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

      document.getElementById("tab"+tabName).click();
      return;
    }
  }
  // If we didn't find a cookie, open the workout tab
  document.getElementById("MedicineLog").click();
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

  // Select a field on the tab based on the name. E.G. For the MedicineLog tab, focus the count input
  if (tabName === "MedicineLog") {
    countInput.focus();
  } else if (tabName === "AvailableMedicines") {
    medicinesTable.rows[medicinesTable.rows.length-1].cells[1].getElementsByTagName('input')[0].focus();
  }

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
