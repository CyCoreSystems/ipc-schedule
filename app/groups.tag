<groups>
   <h4>
      Groups <a class="btn-floating waves-effect waves-light red" onclick={ add } }><i class="material-icons">add</i></a>
   </h3>

   <div class="container">

      <table>
         <thead>
            <th data-field="id">ID</th>
            <th data-field="name">Name</th>
            <th></th>
         </thead>
         <tbody>
            <tr each={opts.groups}>
               <td>{id}</td>
               <td>{name}</td>
               <td>
                  <a class="btn-floating waves-effect waves-light blue" href="#group/{ id }"><i class="material-icons">create</i></a>
                  <a class="btn-floating waves-effect waves-light red" onclick={ remove }><i class="material-icons">delete</i></a>
               </td>
         </tbody>
      </table>

   </div>

   this.add = () => {
      riot.route('/group/new')
   }

   this.remove = () => {
      reload = this.get
      
      fetch('/group/'+opts.group.id,{
         method: 'delete'
      })
      .then(function(resp) {
         if(resp.status == 200) {
            return reload()
         }
      })
   }

   this.get = () => {
      // Fetch the list of groups
      fetch('/groups')
      .then(function(resp) {
         if(resp.status != 200) {
            return console.log("Failed to fetch groups",resp)
         }
         return resp.json()
      })
      .then(function(json) {
         opts.groups = json
         riot.update()
      })
   }

   this.on('mount', function() {
      this.get()
   })
</groups>
	
