var _ = require('lodash');

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
               <td>{ target }</td>
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
         opts.groups = _.sortBy(json,'name')
         self.updateTargets()
         riot.update()
      })
   })

   this.updateTargets = () => {
      opts.groups.forEach(function(g) {
         self.updateTarget(g)
      })
   }

   this.updateTarget = (g) => {
      fetch('/target/'+g.id)
      .then(function(resp) {
         if(resp.status != 200) {
            return console.log("No scheduled target found for group ", g.id)
         }
         resp.text().then(function(t) {
            g.target = t
            var i = _.findIndex(opts.targets, { id: g.id })
            if(i<0) {
               opts.targets.push({id: g.id, target: t})
            } else {
               opts.targets[i].target = t
            }
            riot.update()
         })
      })
   }

</index>
	
