$(document).ready(function () {
    // Show modal to add a new car
    $('#addCarModal').on('show.bs.modal', function () {
        $('#addCarForm')[0].reset(); // Reset form fields
    });

    // Handle form submission for adding a car
    $('#submitCarBtn').on('click', function () {
        var form = $('#addCarForm');
        $.ajax({
            url: '/cars',
            type: 'POST',
            data: form.serialize(),
            success: function () {
                $('#addCarModal').modal('hide');
                location.reload(); // Refresh the page to show the new car
            },
            error: function (xhr, status, error) {
                console.error('Error adding car:', status, error);
                alert('Error adding car. Please check the console for details.');
            }
        });
    });

    // Show modal to edit a car
    $('.editCarBtn').on('click', function () {
        var carId = $(this).data('carid');
        $.get('/cars/' + carId, function (data) {
            $('#editCarID').val(data.ID);
            $('#editNumberPlate').val(data.NumberPlate);
            $('#editMake').val(data.Make);
            $('#editModel').val(data.Model);
            $('#editNumberOfPassengers').val(data.NumberOfPassengers);
            $('#editSaccoID').val(data.SaccoID); // Pre-select the correct SACCO
            $('#editCarModal').modal('show');
        }).fail(function () {
            alert('Error loading car details.');
        });
    });

    // Handle form submission for editing a car
    $('#saveEditCarBtn').on('click', function () {
        var form = $('#editCarForm');
        var carId = $('#editCarID').val();
        $.ajax({
            url: '/cars/edit?id=' + carId,
            type: 'PUT',
            data: form.serialize(),
            success: function () {
                $('#editCarModal').modal('hide');
                location.reload(); // Refresh the page to show the updated car
            },
            error: function () {
                alert('Error updating car.');
            }
        });
    });


    // Delete car
    $('.deleteCarBtn').on('click', function () {
        var carId = $(this).data('carid');
        if (confirm('Are you sure you want to delete this car?')) {
            $.ajax({
                url: '/cars/delete?id=' + carId,
                type: 'DELETE',
                success: function () {
                    location.reload();
                },
                error: function () {
                    alert('Error deleting car.');
                }
            });
        }
    });
});