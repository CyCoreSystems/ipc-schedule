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
         console.log("response:",json)
         opts.groups = json
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
      console.log("Fetching target for group")
      fetch('/target/'+g.id)
      .then(function(resp) {
         console.log("Got response",resp)
         if(resp.status != 200) {
            return console.log("No scheduled target found for group ", g.id)
         }
         resp.text().then(function(t) {
            console.log("Parsed reponse as text",t)
            g.target = t
            var i = _.findIndex(opts.targets, { id: g.id })
            if(i<0) {
               opts.targets.push({id: g.id, target: t})
            } else {
               opts.targets[i].target = t
            }
            console.log("Targets array is now",opts.targets)
            console.log("Groups array is now",opts.groups)
            riot.update()
         })
      })
   }

</index>
	
