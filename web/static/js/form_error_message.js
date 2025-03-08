export async function fetchErrorMessage(errorMessageContainer) {
    try {
        let response =  await fetch('/login');
        if (!response.ok) {
            let data = await response.json().catch(() => null);
            if (data && data.error_message) {
                setErrorMessage(errorMessageContainer, data.error_message);
            }
        }
    } catch (error) {
        console.error(`Error fetching message: ${error}`);
    }
}

function setErrorMessage(errorMessageContainer, message) {
    errorMessageContainer.value = message;
    let errorMessage = document.getElementById('error_text');

    if (errorMessage) {
        errorMessage.textContent = message;
    }
}