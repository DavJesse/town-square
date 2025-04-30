export function renderErrorPage(message, code) {
    // Set page title
    document.title = `Error ${code} - real-time-forum`;

    // Get and clear app container
    const app = document.getElementById('app');
    if (!app) {
        console.error('App container not found');
        return;
    }
    app.innerHTML = '';

    // Create error container
    const errorContainer = document.createElement('section');
    errorContainer.classList.add('error-code');

    // Create and append error title
    const errorTitle = document.createElement('h2');
    errorTitle.textContent = 'Sorry: We Ran Into a Problem';
    errorContainer.appendChild(errorTitle);

    // Create and append error code
    const errorCode = document.createElement('p');
    errorCode.id = 'error_code';
    errorCode.textContent = code;
    errorCode.classList.add('code');
    errorContainer.appendChild(errorCode);

    // Create and append error message
    const errorMessage = document.createElement('p');
    errorMessage.id = 'error_message';
    errorMessage.textContent = message;
    errorMessage.classList.add('issue');
    errorContainer.appendChild(errorMessage);

    // Set error message and code
    setErrorMessage(message, code);

    // Create home link container
    const homeLink = document.createElement('div');
    homeLink.classList.add('home');

    // Create home link
    const homeLinkText = document.createElement('a');
    homeLinkText.id = 'return_home_link';
    homeLinkText.href = '/';
    homeLinkText.textContent = 'Return to Homepage';
    
    // Add click handler to use client-side routing
    homeLinkText.addEventListener('click', (e) => {
        e.preventDefault();
        window.history.pushState({}, '', '/');
        // Assuming you have a router function
        if (typeof navigateTo === 'function') {
            navigateTo('/');
        } else {
            window.location.href = '/';
        }
    });

    homeLink.appendChild(homeLinkText);
    errorContainer.appendChild(homeLink);

    // Add error container to app
    app.appendChild(errorContainer);
}

function setErrorMessage(message, code) {
    const errorMessage = document.getElementById('error_message');
    const errorCode = document.getElementById('error_code');
    
    if (!errorMessage || !errorCode) {
        console.error('Error elements not found');
        return;
    }

    // Handle default 404 case
    if (code === 0 && message === '') {
        errorCode.textContent = '404';
        errorMessage.textContent = 'Not Found';
        return;
    }

    // Map common HTTP status codes to user-friendly messages
    const errorMessages = {
        400: 'Bad Request',
        401: 'Unauthorized',
        403: 'Forbidden',
        404: 'Not Found',
        500: 'Internal Server Error',
        503: 'Service Unavailable'
    };

    errorCode.textContent = code;
    errorMessage.textContent = message || errorMessages[code] || 'Unknown Error';
    
    // Log error for debugging
    console.log(`Error displayed - Code: ${code}, Message: ${message}`);
}

// Export both functions for testing purposes
export { setErrorMessage };
