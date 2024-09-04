$(document).ready(function () {
  // Handle Add Sacco
  $('#submitSaccoBtn').click(function () {
      const saccoName = $('#saccoName').val();
      const manager = $('#manager').val();
      const contact = $('#contact').val();

      if (saccoName && manager && contact) {
          $.post('/home', {
              sacco_name: saccoName,
              manager: manager,
              contact: contact
          }).done(function () {
              location.reload(); // Reload the page to see the new sacco
          }).fail(function () {
              alert("Error adding sacco");
          });
      } else {
          alert("Please fill out all fields");
      }
  });

  // Handle Edit Sacco Button click
  $('.editSaccoBtn').click(function () {
      const saccoID = $(this).data('saccoid');

      // Populate the edit form with current sacco data
      const row = $(this).closest('tr');
      const saccoName = row.find('td:eq(0)').text();
      const manager = row.find('td:eq(1)').text();
      const contact = row.find('td:eq(2)').text();

      $('#editSaccoID').val(saccoID);
      $('#editSaccoName').val(saccoName);
      $('#editManager').val(manager);
      $('#editContact').val(contact);
  });

  // Handle Save Changes for Edit Sacco
  $('#saveEditSaccoBtn').click(function () {
      const saccoID = $('#editSaccoID').val();
      const saccoName = $('#editSaccoName').val();
      const manager = $('#editManager').val();
      const contact = $('#editContact').val();

      if (saccoName && manager && contact) {
          $.post('/edit-sacco', {
              id: saccoID,
              sacco_name: saccoName,
              manager: manager,
              contact: contact
          }).done(function () {
              location.reload(); // Reload the page to see the updated sacco
          }).fail(function () {
              alert("Error editing sacco");
          });
      } else {
          alert("Please fill out all fields");
      }
  });

  // Handle Delete Sacco
  $('.deleteSaccoBtn').click(function () {
      const saccoID = $(this).data('saccoid');
      if (confirm("Are you sure you want to delete this sacco?")) {
          $.post('/delete-sacco', { id: saccoID })
              .done(function () {
                  location.reload(); // Reload the page to remove the deleted sacco
              })
              .fail(function () {
                  alert("Error deleting sacco");
              });
      }
  });
});
