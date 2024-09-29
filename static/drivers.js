$(document).ready(function () {
            // Show modal to add a new driver
            $('#addDriverBtn').on('click', function () {
                $('#addDriverForm')[0].reset();
                $('#addDriverModal').modal('show');
            });

            // Submit driver form
            $('#submitDriverBtn').on('click', function (e) {
                e.preventDefault();
                $.ajax({
                    url: '/drivers',
                    type: 'POST',
                    data: $('#addDriverForm').serializeArray(),
                    success: function () {
                        $('#addDriverModal').modal('hide');
                        location.reload();
                    },
                    error: function (xhr, status, error) {
                        console.error('Error:', status, error);
                        console.error('Response:', xhr.responseText);
                    }
                });
            });

            // Automatically fill car by sacco name
            $(document).ready(function () {
                $('#assignedCar').on('change', function () {
                    var carId = $(this).val();
                    if (carId) {
                        $.get('/api/sacco_by_car', { car_id: carId }, function (data) {
                            $('#saccoID').val(data.ID);
                        });
                    } else {
                        $('#saccoID').val('');
                    }
                });
            });

            // Edit driver
            $('.editDriverBtn').on('click', function () {
                var driverId = $(this).data('driverid');
                $.get('/drivers/' + driverId, function (data) {
                    $('#editDriverID').val(data.ID);
                    $('#editDriverName').val(data.DriverName);
                    $('#editIDNumber').val(data.IDNumber);
                    $('#editContact').val(data.Contact);
                    $('#editAssignedCar').val(data.CarID);
                    $('#editSaccoID').val(data.SaccoID);
                    $('#editDriverModal').modal('show');
                });
            });

            // Save edit driver form
            $('#saveEditDriverBtn').on('click', function () {
                var form = $('#editDriverForm');
                var driverId = $('#editDriverID').val();
                $.ajax({
                    url: '/drivers/edit?id=' + driverId,
                    type: 'PUT',
                    data: $('#editDriverForm').serialize(),
                    success: function () {
                        $('#editDriverModal').modal('hide');
                        location.reload();
                    },
                    error: function (xhr, status, error) {
                        console.error('Error:', status, error);
                        console.error('Response:', xhr.responseText);
                    }
                });
            });
            // Delete driver
            $('.deleteDriverBtn').on('click', function () {
                var driverId = $(this).data('driverid');
                if (confirm('Are you sure you want to delete this driver?')) {
                    $.ajax({
                        url: '/drivers/' + driverId,
                        type: 'DELETE',
                        success: function () {
                            location.reload();
                        }
                    });
                }
            });
        });