export function showError(message) {
    const modalBody = document.getElementById('errorModalBody');
    if (modalBody) {
        modalBody.textContent = message;
        const errorModalElement = document.getElementById('errorModal');
        const errorModal = bootstrap.Modal.getOrCreateInstance(errorModalElement);
        errorModal.show();
    } else {
        alert(message);
    }
}