
<group>
   <div class="container">
      <form id="groupEdit" class="col s12">
         <div class="row s12">
            <input name="id" type="hidden" value={ opts.item.id }/>
            <div class="input-field col s5">
               <input name="name" type="text" class="validate" required value={ opts.item.name } minlength=3 maxlength=25 length=25/>
               <label for="name">Group Name</label>
            </div>
            <div class="input-field col s4">
               <select class="validate" required name="timezone">
                  <option value="" disabled selected={ selected('') }>Choose the timezone for the group</option>
                  <option each={ timezones } value={ value } selected={ selected(value) }>{ label }</option>
               </select>
               <label for="timezone">Time Zone</label>
            </div>
            <div class="input-field col s3">
               <a class="btn-floating waves-effect waves-light green"><i class="material-icons" onclick={ save }>done</i></a>
               <a class="btn-floating waves-effect waves-light red"><i class="material-icons" onclick={ cancel }>cancel</i></a>
            </div>
         </div>
      </form>
   </div>

   this.timezones = [
      {label: 'Eastern', value: 'US/Eastern'},
      {label: 'Central', value: 'US/Central'},
      {label: 'Mountain', value: 'US/Mountain'},
      {label: 'Pacific', value: 'US/Pacific'},
   ];

   opts.item = {}

   this.on('mount', function() {
      if( opts.groupId ) {
         fetch('/group/'+opts.groupId)
         .then(function(resp) {
            if(resp.status == 200) {
               resp.json().then(function(data) {
                  opts.item = data
                  console.log("Item:",opts.item)
                  riot.update()
               })
            }
         })
      }

      $('input').characterCounter()
      $('select').material_select()

   })

   this.selected = (val) => {
      return opts.item.timezone == val
   }

   this.save = () => {
      fetch('/group',{
         method: 'post',
         body: new FormData(this.groupEdit)
      })
      .then(function(resp) {
         if(resp.status != 200) {
            console.log("Failed to save new group")
            return
         }
         parent.editing = false
         return riot.route('/groups')
      })
      .catch(function(ex) {
         alert(ex)
      })
   }

   this.cancel = () => {
      parent.editing = false
      return riot.route('/groups')
   }

</group>
