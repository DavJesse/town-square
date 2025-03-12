export function renderErrorPage() {
    document.title = 'Error!';

    // Create errorContainer
    let errorContainer = document.createElement('div');
    errorContainer.classList.add('error-code');

    // Create h2 to introduce error
    let errorTitle = document.createElement('h2');
    errorTitle.textContent = 'Sorry: We Ran Into a Problem';
    errorContainer.appendChild(errorTitle);

    // Create error message
   let errorCode = document.createElement('p');
   errorCode.id = 'error_code';
   errorContainer.appendChild(errorCode);

    // Create error message
    let errorMessage = document.createElement('p');
    errorMessage.id = 'error_message';
    errorContainer.appendChild(errorMessage);

    fetchErrorMessage(errorContainer, '/error')

    // Add errorContainer to body
    document.body.appendChild(errorContainer);
}

async function fetchErrorMessage(errorContainer, errorLink) {
    try {
        let response =  await fetch(errorLink);
        if (!response.ok) {
            let data = await response.json().catch(() => null);
            
            setErrorMessage(errorContainer, data.issue, data.code);
        }
    } catch (error) {
        console.error(`Error fetching message: ${error}`);
    }
}

function setErrorMessage(errorMessageContainer, message, code) {
    errorMessageContainer.value = message;
    let errorMessage = document.getElementById('error_message');
    let errorCode = document.getElementById('error_code');

    errorCode.textContent = code
        errorMessage.textContent = message;
}