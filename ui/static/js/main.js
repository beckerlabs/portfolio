var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
    var link = navLinks[i];
    if (link.getAttribute('href') == window.location.pathname) {
        link.classList.add("live");
        break;
    }
}

const showOutlineButton = document.getElementById('showOutline');
const showPostsButton = document.getElementById('showPosts');

const outlineSidebar = document.getElementById('outline-sidebar');
const postsSidebar = document.getElementById('posts-sidebar');
const blogPost = document.getElementById('blog-post');

// Function to toggle visibility of a sidebar
function toggleSidebar(button, sidebar, otherButton, otherSidebar) {
    // Toggle the current sidebar
    const isVisible = sidebar.classList.toggle('visible');
    blogPost.classList.toggle('hidden', isVisible);

    // Update the button text
    button.textContent = isVisible ? `Hide ${button.dataset.label}` : button.dataset.label;

    // Hide the other sidebar and reset its button
    if (isVisible) {
        otherSidebar.classList.remove('visible');
        otherButton.textContent = otherButton.dataset.label;
    }
}

// Add data-label attributes to buttons for dynamic text updates
showOutlineButton.dataset.label = 'Outline';
showPostsButton.dataset.label = 'More posts';

// Event listeners for the buttons
showOutlineButton.addEventListener('click', function () {
    toggleSidebar(showOutlineButton, outlineSidebar, showPostsButton, postsSidebar);
});

showPostsButton.addEventListener('click', function () {
    toggleSidebar(showPostsButton, postsSidebar, showOutlineButton, outlineSidebar);
});