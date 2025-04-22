var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
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

showOutlineButton.addEventListener('click', function() {
	outlineSidebar.classList.toggle('visible');
	blogPost.classList.toggle('hidden');
	
	if (outlineSidebar.classList.contains('visible')) {
		showOutlineButton.textContent = 'Hide Outline';
	} else {
		showOutlineButton.textContent = 'Outline';
	}
});

outlineSidebar.addEventListener('click', function(event) {
    if (event.target !== outlineSidebar) {
        outlineSidebar.classList.toggle('visible');
        blogPost.classList.toggle('hidden');
        showOutlineButton.textContent = outlineSidebar.classList.contains('visible') ? 'Hide Outline' : 'Outline';
    }
});

showPostsButton.addEventListener('click', function() {
	postsSidebar.classList.toggle('visible');
	blogPost.classList.toggle('hidden');

	if (postsSidebar.classList.contains('visible')) {
		showPostsButton.textContent = 'Hide Posts';
	} else {
		showPostsButton.textContent = 'More posts';
	}
});