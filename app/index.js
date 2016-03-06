var _ = require('lodash');

// Top menu
require('./menu.tag');
riot.mount("menu")

// Subpages
require('./index.tag');
require('./usage.tag');
require('./upload.tag');
require('./group.tag');
require('./groups.tag');
	

// Routes
riot.route("/", function() {
   riot.mount("#app", "index")
})
riot.route("/group/*", function(id) {
   riot.mount("#app", "group", { groupId: id })
})
riot.route("/groups", function() {
   riot.mount("#app", "groups")
})
riot.route("/usage", function() {
   riot.mount("#app", "usage")
})
riot.route("/upload", function() {
   riot.mount("#app", "upload")
})

// Start the router
riot.route.start(true)
