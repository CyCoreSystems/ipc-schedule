
<upload>

	<div>{status}</div>

	<h4>Upload Weekday ("Day") CSV</h4>
	<form id="dayCsv" onsubmit={uploadDay} method="POST" enctype="multipart/form-data">
		<input id="file" name="file" type="file"/>
		<button type="submit">Upload</button>
	</form>

	<h4>Upload Precise Date ("Date") CSV</h4>
	<form id="dateCsv" onsubmit={uploadDate} method="POST" enctype="multipart/form-data">
		<input id="file" name="file" type="file"/>
		<button type="submit">Upload</button>
	</form>

	<script>
		var self = this;
		self.status = "";
      this.uploadDay = (e) => {
			var form = document.getElementById("dayCsv")
			window.fetch("/sched/import/days", {
				method: "post",
				body: new FormData(form),
			}).then(function(resp) {
				self.status = "Upload Results: " + resp.statusText;
				self.update()
			})
			return false;
		}

      this.uploadDate = (e) => {
			var form = document.getElementById("dateCsv")
			window.fetch("/sched/import/dates", {
				method: "post",
				body: new FormData(form),
			}).then(function(resp) {
				self.status = "Upload Results: " + resp.statusText;
				self.update()
			})
			return false;
		}
	
	</script>

</upload>
