window.renderNavbar = function(data) {
    const navbarHtml = `
            <a href="/">
                <h1 class="navbar__title">forum</h1>
            </a>
            
            <div class="navbar__actions">
                <form id="search-form" class="navbar__search-form">
                    <input type="text" id="search-query" name="q" class="navbar__search" placeholder="Search...">
                    <button type="button" class="navbar__button" onclick="searchPosts()">Search</button>
                </form>
                
                <button class="navbar__button" id="create-post-button">Create Post</button>

                    <form action="/logout" method="POST" style="display: inline;">
                        <button class="navbar__button" id="login-button">Log Out</button>
                    </form>
                    <div class="navbar__Profile__Pic">
                        ${data.ProfPic ? `
                            <button onclick="profile">
                                <img src="/static/images/${data.ProfPic}" alt="user">
                            </button>
                        ` : `
                            <a href="#" onclick="profile()">
                                <img src="/static/user-circle-svgrepo-com.svg" alt="user">
                            </a>
                        `}
                    </div>
            </div>
    `;
    document.querySelector('.navbar').outerHTML = navbarHtml;
    document.getElementById('create-post-button').addEventListener('click', renderCreatePostForm);
}