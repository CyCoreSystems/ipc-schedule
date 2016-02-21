
<index>
	<h4>Schedule<h4>
         <label for="group">Group</label>
         <select name="group" each={ groups }>
            <option id={id} selected={parent.currentGroup}>{name}</option>
         </select>
         <i class="fa fa-plus" onClick={ addGroup }></i>
         <add-group></add-group>


   this.groups = {}

   // Fetch the list of groups
   fetch('/groups')
   .then(function(resp) {
      if(resp.status != 200) {
         return console.log("Failed to fetch groups",resp);
      }
      this.groups = JSON.parse(resp.blob());
   })

</index>
	
