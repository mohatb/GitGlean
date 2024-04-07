document.addEventListener('DOMContentLoaded', function() {
    let allIssues = [];
    let filteredIssues = [];
    let isSearchActive = false;
    let currentPage = 1;
    const issuesPerPage = 10;
    const loader = document.getElementById('loader');


    function renderIssues(issues, page) {
        const start = (page - 1) * issuesPerPage;
        const end = start + issuesPerPage;
        const paginatedIssues = issues.slice(start, end);

        const issuesTableBody = document.getElementById('issuesTable').getElementsByTagName('tbody')[0];
        issuesTableBody.innerHTML = '';

        paginatedIssues.forEach(issue => {
            let row = issuesTableBody.insertRow();
            row.insertCell().textContent = issue.title;
            row.insertCell().textContent = issue.state;

            let link = document.createElement('a');
            link.href = issue.html_url;
            link.textContent = 'View on GitHub';
            link.target = '_blank';
            row.insertCell().appendChild(link);

            row.insertCell().textContent = issue.comments;
            row.insertCell().textContent = issue.reactions.total_count;

            let analyzeButton = document.createElement('button');
            analyzeButton.textContent = 'Analyze';
            analyzeButton.className = 'analyze-btn';
            analyzeButton.onclick = function() { analyzeIssue(issue.url); };
            row.insertCell().appendChild(analyzeButton);
        });

        document.getElementById('page-info').textContent = `Page ${page} of ${Math.ceil(issues.length / issuesPerPage)}`;
    }

    function toggleLoader(show) {
        // Show or hide the loader
        loader.style.display = show ? 'block' : 'none';
    }

    window.changePage = function(delta) {
        const newPage = currentPage + delta;
        const totalIssues = isSearchActive ? filteredIssues : allIssues;
        const totalPages = Math.ceil(totalIssues.length / issuesPerPage);

        if (newPage > 0 && newPage <= totalPages) {
            currentPage = newPage;
            renderIssues(totalIssues, currentPage);
        }
    };

    function performSearch() {
        const searchText = document.getElementById('search-input').value.toLowerCase();
        isSearchActive = searchText.length > 0;
        filteredIssues = allIssues.filter(issue => 
            issue.title.toLowerCase().includes(searchText) || 
            issue.body.toLowerCase().includes(searchText)
        );
        currentPage = 1;
        renderIssues(isSearchActive ? filteredIssues : allIssues, currentPage);
    }

    document.getElementById('search-btn').addEventListener('click', performSearch);
    document.getElementById('search-input').addEventListener('keypress', function(event) {
        if (event.key === 'Enter') {
            event.preventDefault();
            performSearch();
        }
    });

    fetch('/api/list-issues?per_page=30')
        .then(response => response.json())
        .then(data => {
            allIssues = data;
            renderIssues(allIssues, currentPage);
        })
        .catch(error => console.error('Error fetching GitHub issues:', error));

        function analyzeIssue(issueUrl) {
            const issueNumber = issueUrl.split('/').pop();
            toggleLoader(true); // Show loader
            fetch(`/api/analyze-issue?issue_number=${issueNumber}`)
            .then(response => {
                toggleLoader(false); // Hide loader
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                showModal(data.summary);
            })
            .catch(error => {
                console.error('Error analyzing GitHub issue:', error);
                toggleLoader(false); // Hide loader
            });
        }

        function showModal(summary) {
            const modal = document.getElementById('analysisModal');
            const closeSpan = modal.getElementsByClassName('close')[0];
            const analysisText = document.getElementById('analysisText');
            
            // Parse markdown content to HTML
            const parsedContent = marked.parse(summary);
            
            // Set the parsed HTML content
            analysisText.innerHTML = parsedContent;
        
            // Display the modal
            modal.style.display = "block";
        
            // When the user clicks on <span> (x), close the modal
            closeSpan.onclick = function() {
                modal.style.display = "none";
            }
        
            // Close the modal if the user clicks outside of it
            window.onclick = function(event) {
                if (event.target == modal) {
                    modal.style.display = "none";
                }
            }
        }
        

    document.getElementById('search-btn').addEventListener('click', performSearch);

    document.getElementById('search-input').addEventListener('keypress', function(event) {
        if (event.key === 'Enter') {
            event.preventDefault();
            performSearch();
        }
    });

    fetch('/api/list-issues?per_page=30')
    .then(response => response.json())
    .then(data => {
        allIssues = data;
        renderIssues(allIssues, currentPage);
    })
    .catch(error => console.error('Error fetching GitHub issues:', error));
});
