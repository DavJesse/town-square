export function renderErrorPage() {
    document.title = 'Error!';

    // Retriev app from DOM
    let app = document.getElementById('app');
    app.innerHTML = "";

    // Create errorContainer
    let errorContainer = document.createElement('section');
    errorContainer.classList.add('error-code');

    // Create h2 to introduce error
    let errorTitle = document.createElement('h2');
    errorTitle.textContent = 'Sorry: We Ran Into a Problem';
    errorContainer.appendChild(errorTitle);

    // Create error message
   let errorCode = document.createElement('p');
   errorCode.id = 'error_code';
   errorCode.classList.add('code');
   errorContainer.appendChild(errorCode);

    // Create error message
    let errorMessage = document.createElement('p');
    errorMessage.id = 'error_message';
    errorMessage.classList.add('issue');
    errorContainer.appendChild(errorMessage);

    fetchErrorMessage()

    // Create home link container
    let homeLink = document.createElement('div');
    homeLink.classList.add('home');

    // Create home link text
    let homeLinkText = document.createElement('a');
    homeLinkText.id = 'return_home_link';
    homeLinkText.href = '/';
    homeLinkText.textContent = 'Return to Homepage';
    homeLink.appendChild(homeLinkText);

    errorContainer.appendChild(homeLink);

    // Add errorContainer to app
    app.appendChild(errorContainer);
}

async function fetchErrorMessage() {
    try {
        // Fetch the correct error message from `/error`
        let response = await fetch("/error", {
            headers: { "Accept": "application/json" } // Tell Go to return JSON
        });

        let data = await response.json();
        
        // Update the error message
        setErrorMessage(data.issue, data.code, data.path);
    } catch (error) {
        console.error(`Error fetching message: ${error}`);
    }
}


function setErrorMessage(message, code) {
    let errorMessage = document.getElementById('error_message');
    let errorCode = document.getElementById('error_code');
    
    if (code === 0 && message === "") {
        errorCode.textContent = "404";
        errorMessage.textContent = "Not Found";
    } else {
        errorCode.textContent = code;
        errorMessage.textContent = message;
    }
}