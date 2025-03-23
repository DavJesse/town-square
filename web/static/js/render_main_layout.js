
window.renderMainLayout = function() {
    // clear the document body content
    document.body.innerHTML = '';

    // Create navbar
    const navbar = document.createElement('div');
    navbar.classList.add('navbar');
    document.body.appendChild(navbar);

    // Create main container
    const container = document.createElement('div');
    container.classList.add('container');
    document.body.appendChild(container);

    // Create hamburger button for mobile devices
    const hamburgerButton = document.createElement('button');
    hamburgerButton.classList.add('hamburger');
    hamburgerButton.onclick = function() { toggleSidebar(); };
    container.appendChild(hamburgerButton);

    // Create sidebar
    const sidebar = document.createElement('div');
    sidebar.classList.add('sidebar');
    sidebar.id = 'sidebar';
    container.appendChild(sidebar);

    // Create close button for sidebar
    const closeButton = document.createElement('button');
    closeButton.classList.add('sidebar__close');
    closeButton.onclick = function() { toggleSidebar(); };
    sidebar.appendChild(closeButton);

    // Create sidebar heading
    const sidebarHeading = document.createElement('h3');
    sidebarHeading.classList.add('sidebar__heading');
    sidebarHeading.textContent = 'Categories';
    sidebar.appendChild(sidebarHeading);

    // Create categories list container
    const categoriesList = document.createElement('div');
    categoriesList.id = 'categories-list';
    sidebar.appendChild(categoriesList);

    // Create main content area
    const mainContent = document.createElement('div');
    mainContent.classList.add('main-content');
    mainContent.id = 'main-content';
    container.appendChild(mainContent);
}
