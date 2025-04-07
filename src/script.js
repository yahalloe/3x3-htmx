function scrollToTop() {
  window.scrollTo({
    top: 0,
    behavior: "smooth"
  });
}


// This will work for both initial page load and HTMX content swaps
document.addEventListener("DOMContentLoaded", setupLoadingScreen);
document.addEventListener("htmx:afterSwap", setupLoadingScreen);

function setupLoadingScreen() {
  // Hide loading screen by default
  const loadingScreen = document.getElementById("loading-screen");
  if (loadingScreen) loadingScreen.style.display = "none";

  // Event delegation for external links
  document.body.addEventListener("click", function(e) {
    const link = e.target.closest("a[href^='https']");
    
    if (link && !link.hasAttribute("hx-get")) {
      e.preventDefault();
      loadingScreen.style.display = "flex";
      
      setTimeout(() => {
        window.location.href = link.href;
      }, 300);
    }
  });
}

// Additional safety measure for back button
window.addEventListener('pageshow', function(event) {
  if (event.persisted) {
    const loadingScreen = document.getElementById("loading-screen");
    if (loadingScreen) loadingScreen.style.display = "none";
  }
});