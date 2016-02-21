
<group>

   <div if={ !editing }>
      {name}<a class="btn-floating waves-effect waves-light red"><i class="material-icons" onclick={ edit }>settings</i></a>
   </div>
   <form if={editing} class="col s12">
      <div class="row">
         <div class="input-field col s6">
            <input placeholder="Group Name" name="name" type="text" class={valid ? "validate"}/>
            <label for="name">Group Name</label>
         </div>
         <div class="input-field col s5">
            <select name="timezone">
               <option each={ timezones } value={ value } selected={ selected }>{ label }</option>
            </select>
            <label for="timezone">Time Zone</label>
         </div>
         <div class="input-field col s5">
            <select name="timezone">
               <option each={ timezones } value={ value } selected={ selected }>{ label }</option>
            </select>
            <label for="timezone">Time Zone</label>
         </div>
      </div>
   </form>

   this.timezones = [
      {label: 'Eastern', value: 'US/Eastern'},
      {label: 'Central', value: 'US/Central'},
      {label: 'Mountain', value: 'US/Mountain'},
      {label: 'Pacific', value: 'US/Pacific'},
   ];

   this.editing = false;

   edit() {
      this.editing = true;
   }

</group>
	
