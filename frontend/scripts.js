document.getElementById('upload_button').addEventListener('click', function () {
    var fileInput = document.getElementById('formFile');
    if (!fileInput || !fileInput.files || fileInput.files.length === 0) {
        Swal.fire({
            icon: 'error',
            title: 'No image selected',
            text: 'Please choose an image to upload.',
            showConfirmButton: false,
            timerProgressBar: true,
            timer: 1500
        });
        return;
    }

    Swal.fire({
        title: 'Are you sure?',
        text: 'Do you want to proceed with the upload?',
        icon: 'question',
        showCancelButton: true,
        confirmButtonText: 'Yes, proceed',
        cancelButtonText: 'No, cancel',
    }).then((result) => {
        if (result.isConfirmed) {
            document.getElementById('uploading_form').submit();
        } else {
            Swal.fire({
                icon: 'info',
                title: 'Upload canceled',
                showConfirmButton: false,
                timerProgressBar: true,
                timer: 1500
            });
        }
    });
});
