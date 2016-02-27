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
   
   var self = this

   opts.groups = []
   opts.targets = []

   this.on('mount', function() {
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
         self.updateTargets()
         riot.update()
      })
   })

   updateTargets() {
      opts.groups.forEach(function(g) {
         self.updateTarget(g)
      })
   }

   updateTarget(g) {
      fetch('/target/'+g.id)
      .then(function(resp) {
         if(resp.status != 200) {
            return console.log("No schedule found for group ", g.id)
         }
         var i = _.findIndex(opts.targets, { id: g.id })
         resp.text().then(function(t) {
            if(i<0) {
               opts.targets.append({id: g.id, target: t})
            } else {
               opts.targets[i].target = t
            }
            self.update()
         })
      })
   }

   getTarget() {
      t = _.find(opts.targets,{id: item.id})
      return typeof t !== undefined ? t.target : ''
   }
</index>
	
