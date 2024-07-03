package provider

import (
	"context"
	"fmt"
	"strconv"

	"github.com/SeoGeonhyuk/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &gameResource{}
	_ resource.ResourceWithConfigure   = &gameResource{}
	_ resource.ResourceWithImportState = &gameResource{}
)

// NewOrderResource is a helper function to simplify the provider implementation.
func NewGameResource() resource.Resource {
	return &gameResource{}
}

// orderResource is the resource implementation.
type gameResource struct {
	client *hashicups.Client
}

// ImportState implements resource.ResourceWithImportState.
func (r *gameResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// gameResourceModel maps the resource schema data.
type gameResourceModel struct {
	ID          types.String     `tfsdk:"id"`
	Name      	types.String 	 `tfsdk:"name"`
	StarPoint   types.Float64    `tfsdk:"star_point"`
	PlayerNum   types.Int64		 `tfsdk:"player_num"`
}


// Configure implements resource.ResourceWithConfigure.
func (r *gameResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*hashicups.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *hashicups.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.client = client
}

// Metadata returns the resource type name.
func (r *gameResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_game"
}

// Schema defines the schema for the resource.
func (r *gameResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"star_point": schema.Float64Attribute{
				Required: true,
			},
			"player_num": schema.Int64Attribute{
				Required: true,
			},
		},
	}
}

// Create a new resource.
func (r *gameResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan gameResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API request body from plan
	// var items []hashicups.OrderItem
	// for _, item := range plan. {
	// 	items = append(items, hashicups.OrderItem{
	// 		Coffee: hashicups.Coffee{
	// 			ID: int(item.Coffee.ID.ValueInt64()),
	// 		},
	// 		Quantity: int(item.Quantity.ValueInt64()),
	// 	})
	// }

	newGame := hashicups.Game{
		Name:     plan.Name.ValueString(),
        StarPoint: float32(plan.StarPoint.ValueFloat64()),
        PlayerNum: int(plan.PlayerNum.ValueInt64()),
    }
	// Create new order
	game, err := r.client.CreateGame(newGame, &r.client.Token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating order",
			"Could not create order, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	plan.ID = types.StringValue(strconv.Itoa(game.ID))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information.
func (r *gameResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state gameResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get refreshed order value from HashiCups
	game, err := r.client.GetGame(state.ID.ValueString(), &r.client.Token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading HashiCups Order",
			"Could not read HashiCups order ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}



	// Set refreshed state
	diags = resp.State.Set(ctx, &game)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *gameResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan gameResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	updateGame := hashicups.Game{
		Name:     plan.Name.ValueString(),
        StarPoint: float32(plan.StarPoint.ValueFloat64()),
        PlayerNum: int(plan.PlayerNum.ValueInt64()),
    }

	// Update existing order
	_, err := r.client.UpdateGame(plan.ID.ValueString(), updateGame, &r.client.Token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating HashiCups Order",
			"Could not update order, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated items from GetOrder as UpdateOrder items are not
	// populated.
	game, err := r.client.GetGame(plan.ID.ValueString(), &r.client.Token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading HashiCups Order",
			"Could not read HashiCups order ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	

	diags = resp.State.Set(ctx, game)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *gameResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state gameResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete existing order
	err := r.client.DeleteGame(state.ID.ValueString(), &r.client.Token)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting HashiCups Order",
			"Could not delete order, unexpected error: "+err.Error(),
		)
		return
	}
}
