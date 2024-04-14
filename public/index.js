const allPostDates = document.querySelectorAll("div.post-time");
if (allPostDates.length) {
    for (const postDateDiv of allPostDates) {
        const text = postDateDiv.textContent.trim();
        const formattedDate = new Date(text);
        postDateDiv.textContent = formattedDate.toDateString();
    }
}