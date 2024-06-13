package provider

import (
	"context"
	"fmt"
	"github.com/bunnyway/terraform-provider-bunny/internal/api"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
)

var _ resource.Resource = &StorageResource{}
var _ resource.ResourceWithImportState = &StorageResource{}

func NewStorageResource() resource.Resource {
	return &StorageResource{}
}

type StorageResource struct {
	client *api.Client
}

type StorageResourceModel struct {
	Id                 types.Int64  `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Password           types.String `tfsdk:"password"`
	ReadOnlyPassword   types.String `tfsdk:"password_readonly"`
	Region             types.String `tfsdk:"region"`
	ReplicationRegions types.Set    `tfsdk:"replication_regions"`
	StorageHostname    types.String `tfsdk:"hostname"`
	ZoneTier           types.String `tfsdk:"zone_tier"`
	Custom404FilePath  types.String `tfsdk:"custom_404_file_path"`
	Rewrite404To200    types.Bool   `tfsdk:"rewrite_404_to_200"`
	DateModified       types.String `tfsdk:"date_modified"`
}

var storageZoneTierMap = map[uint8]string{
	0: "Standard",
	1: "Edge",
}

func (r *StorageResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storagezone"
}

func (r *StorageResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Storage Zone",

		Attributes: map[string]schema.Attribute{
			"id": schema.Int64Attribute{
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"region": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"replication_regions": schema.SetAttribute{
				ElementType: types.StringType,
				Optional:    true,
			},
			"zone_tier": schema.StringAttribute{
				Required: true,
			},
			"hostname": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"password_readonly": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"custom_404_file_path": schema.StringAttribute{
				Optional: true,
			},
			"rewrite_404_to_200": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"date_modified": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *StorageResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*api.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *api.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *StorageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var dataTf StorageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &dataTf)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dataApi := r.convertModelToApi(ctx, dataTf)
	dataApi, err := r.client.CreateStoragezone(dataApi)
	if err != nil {
		resp.Diagnostics.AddError("Unable to create storage zone", err.Error())
		return
	}

	tflog.Trace(ctx, "created storagezone "+dataApi.Name)
	dataTf, diags := r.convertApiToModel(dataApi)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	if len(dataApi.ReplicationRegions) == 0 {
		dataTf.ReplicationRegions = types.SetNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &dataTf)...)
}

func (r *StorageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data StorageResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	dataApi, err := r.client.GetStoragezone(data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("Error fetching storage zone", err.Error()))
		return
	}

	dataTf, diags := r.convertApiToModel(dataApi)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &dataTf)...)
}

func (r *StorageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data StorageResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// @TODO removing a region from replication_regions is not allowed
	dataApi := r.convertModelToApi(ctx, data)
	dataApi, err := r.client.UpdateStoragezone(dataApi)
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("Error updating storage zone", err.Error()))
		return
	}

	dataTf, diags := r.convertApiToModel(dataApi)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &dataTf)...)
}

func (r *StorageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data StorageResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteStoragezone(data.Id.ValueInt64())
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("Error deleting storage zone", err.Error()))
	}
}

func (r *StorageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := strconv.ParseInt(req.ID, 10, 64)
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("Error converting ID to integer", err.Error()))
		return
	}

	dataApi, err := r.client.GetStoragezone(id)
	if err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("Error fetching storage zone", err.Error()))
		return
	}

	dataTf, diags := r.convertApiToModel(dataApi)
	if diags != nil {
		resp.Diagnostics.Append(diags...)
		return
	}

	if len(dataApi.ReplicationRegions) == 0 {
		dataTf.ReplicationRegions = types.SetNull(types.StringType)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &dataTf)...)
}

func (r *StorageResource) convertModelToApi(ctx context.Context, dataTf StorageResourceModel) api.Storagezone {
	dataApi := api.Storagezone{}
	dataApi.Id = dataTf.Id.ValueInt64()
	dataApi.Name = dataTf.Name.ValueString()
	dataApi.Password = dataTf.Password.ValueString()
	dataApi.ReadOnlyPassword = dataTf.ReadOnlyPassword.ValueString()
	dataApi.ZoneTier = mapValueToKey(storageZoneTierMap, dataTf.ZoneTier.ValueString())
	dataApi.Region = dataTf.Region.ValueString()
	dataApi.Rewrite404To200 = dataTf.Rewrite404To200.ValueBool()
	dataApi.Custom404FilePath = dataTf.Custom404FilePath.ValueString()
	dataApi.StorageHostname = dataTf.StorageHostname.ValueString()
	dataApi.DateModified = dataTf.DateModified.ValueString()

	if dataTf.ReplicationRegions.IsNull() {
		dataApi.ReplicationRegions = nil
	} else {
		replicationRegions := []string{}
		dataTf.ReplicationRegions.ElementsAs(ctx, &replicationRegions, false)
		dataApi.ReplicationRegions = replicationRegions
	}

	return dataApi
}

func (r *StorageResource) convertApiToModel(dataApi api.Storagezone) (StorageResourceModel, diag.Diagnostics) {
	dataTf := StorageResourceModel{}
	dataTf.Id = types.Int64Value(dataApi.Id)
	dataTf.Name = types.StringValue(dataApi.Name)
	dataTf.Password = types.StringValue(dataApi.Password)
	dataTf.ReadOnlyPassword = types.StringValue(dataApi.ReadOnlyPassword)
	dataTf.ZoneTier = types.StringValue(mapKeyToValue(storageZoneTierMap, dataApi.ZoneTier))
	dataTf.Region = types.StringValue(dataApi.Region)
	dataTf.Rewrite404To200 = types.BoolValue(dataApi.Rewrite404To200)
	dataTf.StorageHostname = types.StringValue(dataApi.StorageHostname)
	dataTf.DateModified = types.StringValue(dataApi.DateModified)

	if len(dataApi.ReplicationRegions) == 0 {
		dataTf.ReplicationRegions = types.SetNull(types.StringType)
	} else {
		regions := make([]attr.Value, len(dataApi.ReplicationRegions))
		for i, region := range dataApi.ReplicationRegions {
			regions[i] = types.StringValue(region)
		}

		replicationRegions, err := types.SetValue(types.StringType, regions)
		if err != nil {
			return dataTf, err
		}

		dataTf.ReplicationRegions = replicationRegions
	}

	if len(dataApi.Custom404FilePath) > 0 {
		dataTf.Custom404FilePath = types.StringValue(dataApi.Custom404FilePath)
	}

	return dataTf, nil
}