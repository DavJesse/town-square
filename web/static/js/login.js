import { fetchErrorMessage } from '/static/js/form_error_message.js'

export function renderLoginPage() {
    // update document title
    document.title = 'login - real-time-forum';

    // create login container
    let loginContainer = document.createElement('div');
    loginContainer.classList.add('login-container');
    
    // create logo
    let heading = document.createElement('h1');
    heading.classList.add('logo');
    heading.textContent = 'real-time-forum';
    loginContainer.appendChild(heading);

    // create login sub-container
    let loginSubContainer = document.createElement('div');
    loginSubContainer.classList.add('login');

    // create login tagline
    let loginTagline = document.createElement('h2');
    loginTagline.classList.add('login-tagline');
    loginTagline.textContent = 'Welcome back, sign in to continue';
    loginSubContainer.appendChild(loginTagline);

    // create error message container
    let errorMessageContainer = document.createElement('input');
    errorMessageContainer.classList.add('error-message');
    errorMessageContainer.id = 'error_message';
    errorMessageContainer.type = 'hidden';
    errorMessageContainer.value = '';

    // Create error message placeholder
    let errorText = document.createElement('p');
    errorText.id = 'error_text';
    errorText.style.color = 'red';
    loginSubContainer.appendChild(errorText);
    loginSubContainer.appendChild(errorMessageContainer);
    
    fetchErrorMessage(errorMessageContainer);

    // create login form
    let loginForm = document.createElement('form');
    loginForm.method = 'POST';
    loginForm.action = '/login';

    // create email/username input
    let emailUsername = document.createElement('input');
    emailUsername.type = 'text';
    emailUsername.id = 'email_username';
    emailUsername.placeholder = 'username/email';
    emailUsername.name = 'email_username';
    emailUsername.required = true;
    loginForm.appendChild(emailUsername);

    // create password input
    let password = document.createElement('input');
    password.type = 'password';
    password.id = 'password';
    password.placeholder = 'password';
    password.name = 'password';
    password.required = true;
    loginForm.appendChild(password)
    
    // create submit button
    let loginButton = document.createElement('button');
    loginButton.type ='submit';
    loginButton.textContent = 'login';
    loginForm.appendChild(loginButton)
    
    loginSubContainer.appendChild(loginForm);
    loginContainer.appendChild(loginSubContainer); // add login sub-container

    // create account issues
    let accountIssues = document.createElement('div');
    accountIssues.classList.add('account-issues');

    let registerOption = document.createElement('p');
    registerOption.textContent = 'Don\'t have an account? ';

    let registerLink = document.createElement('a');
    registerLink.classList.add('register-link');
    registerLink.textContent = 'Sign up';
    registerLink.href = '/register';

    registerOption.appendChild(registerLink);
    accountIssues.appendChild(registerOption);
    loginContainer.appendChild(accountIssues);

    document.body.appendChild(loginContainer);
    
    

        // Attach event listener to handle login via AJAX
        loginForm.addEventListener("submit", async function (event) {
            event.preventDefault(); // Prevent full-page reload
        
            if (event.target !== loginForm) return; // Ensure accurate form submision

            let emailUsername = document.getElementById("email_username").value;
            let password = document.getElementById("password").value;
            
            let requestBody = JSON.stringify({
                email_username: emailUsername,
                password: password,
            });
            
            try {
            let response = await fetch("/login", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: requestBody,
            });
        
            let data = await response.json().catch(() => null);
            
            if (response.ok && !data?.error_message) {
                // Redirect to dashboard on success
                navigateTo("/");
            } else {
                // Show error message
                document.getElementById("error_text").textContent = data?.error_message || "Login failed";
            }
        } catch (error) {
            console.log(`Fetch Error: ${error}`);
            document.getElementById("error_text").textContent = "Network Error, please try again";
        }
    });
    
    // Prevent 'sign up' link from being blocked
    registerLink.addEventListener("click", function(event) {
        event.stopPropagation();
        navigateTo("/register");
    });
}

