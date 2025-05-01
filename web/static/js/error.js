export function renderErrorPage(message, code) {
    message = message ?? 'Not Found';
    code = code ?? 404;

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

    // Create home link container
    const homeLink = document.createElement('div');
    homeLink.classList.add('home');

    // Create home link
    const homeLinkText = document.createElement('a');
    homeLinkText.id = 'return_home_link';
    homeLinkText.href = '/';
    homeLinkText.textContent = 'Return to Homepage';

    homeLink.appendChild(homeLinkText);
    errorContainer.appendChild(homeLink);

    // Add error container to app
    app.appendChild(errorContainer);
}

