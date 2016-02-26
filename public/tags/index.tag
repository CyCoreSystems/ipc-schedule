<index>
   <h4>Schedule</h4>

   <div class="container">
      <table>
         <thead>
            <tr>
               <th>Group</th>
               <th>Target</th>
            </tr>
         </thead>
         <tbody>
            <tr each={ opts.groups }>
               <td>{ name }</td>
               <td>{ getTarget(id) }</td>
            </tr>
         </tbody>
      </table>
   </div>
   
   opts.groups = []

   // Fetch the list of groups
   fetch('/groups')
   .then(function(resp) {
      if(resp.status != 200) {
         return console.log("Failed to fetch groups",resp)
      }
      return resp.json()
   })
   .then(function(json) {
      console.log("response:",json)
      opts.groups = json
      riot.update()
   })

   getTarget(id) {
   }
</index>
	
