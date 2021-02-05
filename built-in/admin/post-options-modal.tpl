<div class="modal-header">
    <h3 class="modal-title">文章选项</h3>
</div>
<div class="modal-body">
    <div class="container-fluid">
        <form class="form-horizontal">
            <div class="form-group">
                <label for="post-slug" class="col-sm-2 control-label">自定义Slug</label>
                <div class="col-sm-4">
                    <input spellcheck="true" type="text" class="form-control" id="post-slug" ng-model="shared.post.Slug" value="{{shared.post.Slug}}">
                </div>
            </div>
            <div class="form-group">
                <label for="post-meta-description" class="col-sm-2 control-label">文章简述</label>
                <div class="col-sm-4">
                    <input spellcheck="true" type="text" class="form-control" id="post-meta-description" ng-model="shared.post.MetaDescription" value="{{shared.post.MetaDescription}}">
                </div>
            </div>
            <div class="form-group">
                <label for="post-cover" class="col-sm-2 control-label">封面</label>
                <div class="col-sm-10">
                    <a ng-controller="ImageModalCtrl" ng-click="open('lg', 'post-cover')"><img class="img-settings img-thumbnail img-settings" id="post-cover" src="{{shared.post.Image}}" alt="{{shared.post.Image}}" ng-if="shared.post.Image!=''" /><img class="img-settings img-thumbnail img-settings" id="post-cover" src="/public/images/no-image.png" alt="No image" ng-if="shared.post.Image==''" /></a> <a class="text-danger" id="post-cover-delete" ng-controller="EmptyModalCtrl" ng-click="deleteCover()"><span class="glyphicon glyphicon-remove" aria-hidden="true"></span> 删除</a>
                </div>
            </div>
            <div class="form-group">
                <div class="col-sm-offset-1 col-sm-10">
                  <div class="checkbox">
                    <label>
                        <input bs-switch ng-model="shared.post.IsFeatured" type="checkbox" class="post-checkbox" data-label-text="Feature Post" data-label-width="85" data-off-text="否" data-on-text="是" data-on-color="success" data-off-color="danger" data-size="normal">
                    </label>
                  </div>
                </div>
            </div>
            <div class="form-group">
                <div class="col-sm-offset-1 col-sm-10">
                  <div class="checkbox">
                    <label>
                        <input bs-switch ng-model="shared.post.IsPage" type="checkbox" class="post-checkbox" data-label-text="静态页面" data-label-width="85" data-off-text="否" data-on-text="是" data-on-color="success" data-off-color="danger" data-size="normal">
                    </label>
                  </div>
                </div>
            </div>
        </form>
    </div>
</div>
<div class="modal-footer">
    <button class="btn btn-primary" ng-click="ok()">确定</button>
</div>