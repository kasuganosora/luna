<div class="modal-header">
    <h3 class="modal-title">上传附件</h3>
</div>
<div class="modal-body">
	<div class="container-fluid">
		<input id="file-input" name="multiplefiles" type="file" multiple=true class="file-loading" data-upload-url="/admin/api/upload" data-max-file-count="10">
		<div class="modal-footer modal-divider">
			<button class="btn btn-primary" ng-click="ok()">确定</button>
			<button class="btn btn-warning" ng-click="cancel()">取消</button>
		</div>
		<div infinite-scroll="shared.infiniteScrollFactory.nextPage()" infinite-scroll-disabled="infiniteScrollFactory.busy" infinite-scroll-distance="1" infinite-scroll-container="'.modal'">
			<div class="col-xs-4 col-sm-3" ng-if="$index<shared.infiniteScrollFactory.items.length" ng-repeat="image in shared.infiniteScrollFactory.items track by $index">
				<a class="instance-hook" ng-click="shared.selected = shared.infiniteScrollFactory.items[$index]" img-selection-directive>
					<img ng-class="{imgselected:$first,firstimg:$first}" class="img-thumbnail img-modal center-block" ng-src="{{shared.infiniteScrollFactory.items[$index]}}" alt="{{shared.infiniteScrollFactory.items[$index]}}" />
				</a>
				<div id="image-delete"><a class="text-danger" ng-click="deleteImage(shared.infiniteScrollFactory.items[$index])"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span> 删除</a></div>
			</div>
		</div>
	</div>
</div>
<div class="modal-footer">
	<button class="btn btn-primary" ng-click="ok()">确定</button>
	<button class="btn btn-warning" ng-click="cancel()">取消</button>
</div>
<script>
	$("#file-input").fileinput({language:"zh"});
	$('#file-input').on('fileuploaded', function(event, data, previewId, index) {
		$('#file-input').fileinput('reset');
		angular.element($("[ng-controller='ImageModalCtrl']")).scope().$apply(function () {
			var infiniteScrollFactory = angular.element($("[ng-controller='ImageModalCtrl']")).scope().shared.infiniteScrollFactory;
			infiniteScrollFactory.after = 1;
			infiniteScrollFactory.busy = false;
			infiniteScrollFactory.items = [];
			infiniteScrollFactory.nextPage();
			angular.element($("[ng-controller='ImageModalCtrl']")).scope().shared.selected = data.response[0];
		});
	});
</script>