package repositories

import (
	"DotConnect/schemas"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type WorkspaceRepository interface {
	CreateWorkspace(ctx context.Context, workspacePayload *schemas.WorkSpace) (*schemas.WorkSpace, error)
	GetAllWorkspaces(ctx context.Context, userId string) ([]*schemas.WorkSpace, error)
	GetSlug(ctx context.Context,slug string)(*schemas.WorkSpace, error)
	SlugAlreadyExists(ctx context.Context, slug string) (bool, error)
	AddMemberToWorkspace(ctx context.Context,slug string,member schemas.WorkspaceMember)(bool,error)
	MemberExists(ctx context.Context,slug string,userId string) (bool, error)
	AlreadyAdmin(ctx context.Context,slug string,userId string,role string) (bool, error)
	UpdateMemberRole(ctx context.Context,slug string,targetUser string, newRole string)(bool,error)
	GetAllMembers(ctx context.Context, slug string) ([]schemas.WorkspaceMember, error) 
	CreateInvite(ctx context.Context,slug string,invite schemas.WorkspaceInvite) error
	GetInviteByCode(ctx context.Context,code string)(*schemas.WorkSpace, *schemas.WorkspaceInvite, error)
	CreateJoinRequest(ctx context.Context,slug string,req *schemas.JoinRequest) error
	GetJoinRequest(ctx context.Context,slug string) ([] schemas.JoinRequest,error)
	ApproveRequest(ctx context.Context,slug string,inviteCode string,member schemas.WorkspaceMember,requestID bson.ObjectID,reviewerId string)(string,error)
	RejectRequest(ctx context.Context,slug string,inviteCode string,requestID bson.ObjectID,reviewerId string)(string,error)
	UpdateWorkspace(ctx context.Context,slug string, mData *schemas.WorkSpace) (*schemas.WorkSpace,error)
	RemoveMemberFromWorkspace(ctx context.Context,slug string,memberId string)(bool ,error)
	DeleteWorkspace(ctx context.Context,slug string)(bool,error)
	
}

type WorkspaceRepositoryImpl struct {
	db *mongo.Collection
}

func NewWorkspaceRepository(_db *mongo.Collection) *WorkspaceRepositoryImpl {
	return &WorkspaceRepositoryImpl{
		db: _db,
	}
}

func (repo *WorkspaceRepositoryImpl) CreateWorkspace(ctx context.Context, workspacePayload *schemas.WorkSpace) (*schemas.WorkSpace, error) {

	workspacePayload.ID = bson.NewObjectID()
	workspacePayload.CreatedAt = time.Now()
	workspacePayload.UpdatedAt = time.Now()

	_, err := repo.db.InsertOne(ctx, workspacePayload)

	if err != nil {
		return nil, fmt.Errorf("error occurred while creating the workspace: %w", err)
	}
	return workspacePayload, nil
}

func (repo *WorkspaceRepositoryImpl) GetAllWorkspaces(ctx context.Context, userId string) ([]*schemas.WorkSpace, error) {

	cursor, err := repo.db.Find(ctx, bson.M{"members.userId": userId}, options.Find().SetSort(bson.D{{Key: "createdAt", Value: -1}}))
	if err != nil {
		return nil, fmt.Errorf("error fetching workspaces: %w", err)
	}
	defer cursor.Close(ctx)

	var workspaces []*schemas.WorkSpace

	err = cursor.All(ctx, &workspaces)
	if err != nil {
		return nil, fmt.Errorf("error decoding workspaces: %w", err)
	}
	return workspaces, nil
}

func (repo *WorkspaceRepositoryImpl) GetAllMembers(ctx context.Context, slug string) ([]schemas.WorkspaceMember, error)  {

	var ws schemas.WorkSpace
	err := repo.db.FindOne(ctx, bson.M{"slug": slug}).Decode(&ws)
	if err != nil {
		return nil, fmt.Errorf("error fetching workspace members: %w", err)
	}
	return ws.Members, nil
}

func (repo *WorkspaceRepositoryImpl) SlugAlreadyExists(ctx context.Context, slug string) (bool, error) {
	count, err := repo.db.CountDocuments(ctx, bson.M{"slug": slug})
	return count > 0, err
}

func (repo *WorkspaceRepositoryImpl)MemberExists(ctx context.Context,slug string,userId string) (bool, error){
	
	count, err := repo.db.CountDocuments(ctx, bson.M{"slug": slug,"members.userId":userId})
	if err != nil {
		return false, err
	}

	return count > 0, err
}

func (repo *WorkspaceRepositoryImpl)AlreadyAdmin(ctx context.Context,slug string,userId string,role string) (bool, error){
	

	filter:=bson.M{
		"slug":slug,
		"members":bson.M{"$elemMatch": bson.M{
				"userId": userId,
				"role":   role,
			},
		},
	}

	count, err := repo.db.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}

	return count > 0, err
}

func(repo *WorkspaceRepositoryImpl) AddMemberToWorkspace(ctx context.Context,slug string,member schemas.WorkspaceMember)(bool,error){
	
	filter:= bson.M{"slug":slug}
	update:=bson.M{
		"$addToSet":bson.M{"members":member},
		"$set":bson.M{"JoinedAt":time.Now()},
	}


	result,err := repo.db.UpdateOne(ctx,filter,update)
	if err != nil{
		return false,fmt.Errorf("error adding member: %w",err)
	}

	if result.MatchedCount == 0{
		return false,fmt.Errorf("workspace not found")
	}

	return true,nil
}

func (repo *WorkspaceRepositoryImpl) GetSlug(ctx context.Context,slug string)(*schemas.WorkSpace, error){

	var workspaces *schemas.WorkSpace
	err := repo.db.FindOne(ctx,bson.M{"slug":slug}).Decode(&workspaces)
	if err != nil{
		return nil, err
	}

	return workspaces,nil
}
func (repo *WorkspaceRepositoryImpl)UpdateMemberRole(ctx context.Context,slug string,targetUser string, newRole string)(bool,error){

	filter:=bson.M{"slug":slug,"members.userId":targetUser}
	update:=bson.M{
		"$set":bson.M{"members.$.role":newRole,
			"updatedAt":time.Now()},
	}

	result,err := repo.db.UpdateOne(ctx,filter,update)
	if err != nil{
		return false,fmt.Errorf("error updating member role: %w",err)
	}

	if result.MatchedCount == 0{
		return false,fmt.Errorf("workspace not found or user is not a member of this workspace")
	}

	return true,nil
}

func (repo *WorkspaceRepositoryImpl)CreateInvite(ctx context.Context,slug string,invite schemas.WorkspaceInvite) error{
	filter:= bson.M{"slug":slug}
	update := bson.M{
		"$push": bson.M{"invites": invite},
		"$set":  bson.M{"updatedAt": time.Now()},
	}
	result, err := repo.db.UpdateOne(ctx, filter, update)
	if err != nil {
		return fmt.Errorf("error creating invite: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("workspace not found")
	}
	return nil
}

func (repo *WorkspaceRepositoryImpl)GetInviteByCode(ctx context.Context,code string)(*schemas.WorkSpace, *schemas.WorkspaceInvite, error){
	filter:=bson.M{"invites.code":code}

	var ws schemas.WorkSpace
	err := repo.db.FindOne(ctx,filter).Decode(&ws)
	if err != nil{
		if err == mongo.ErrNoDocuments {
			return nil, nil, fmt.Errorf("invalid invite link")
		}
		return nil, nil, fmt.Errorf("error finding invite: %w", err)
	}

	for _, inv := range ws.Invites {
		if inv.Code == code {
			return &ws, &inv, nil
		}
	}
	return nil, nil, fmt.Errorf("invite code not found")
}

func (repo *WorkspaceRepositoryImpl) CreateJoinRequest(ctx context.Context,slug string,req *schemas.JoinRequest) error{
	
	filter:= bson.M{"slug":slug}
	update := bson.M{
		"$push": bson.M{"joinRequests": req},
		"$set":  bson.M{"updatedAt": time.Now()},
	}

	result,err := repo.db.UpdateOne(ctx,filter,update)
	if err!= nil{
		return fmt.Errorf("Error occured while submitting the joining request")
	}
	if result.MatchedCount == 0{
		return fmt.Errorf("No workspace found")
	}
	return nil
}

func (repo *WorkspaceRepositoryImpl) GetJoinRequest(ctx context.Context,slug string) ([] schemas.JoinRequest,error){
	
	var ws *schemas.WorkSpace
	filter:=bson.M{"slug":slug}
	err := repo.db.FindOne(ctx,filter).Decode(&ws)
	if err != nil{
		return nil,fmt.Errorf("No workspace found")
	}
	return ws.JoinRequests,nil
}

func (repo *WorkspaceRepositoryImpl)ApproveRequest(ctx context.Context,slug string,inviteCode string,member schemas.WorkspaceMember,requestID bson.ObjectID,reviewerId string)(string,error){

	var ws schemas.WorkSpace
	var invitationLink *schemas.WorkspaceInvite


	filter := bson.M{"slug":slug}

	err:= repo.db.FindOne(ctx,filter).Decode(&ws)
	if err != nil{
		return "",fmt.Errorf("workspace not found %w",err)
	}

	for _,invitation := range ws.Invites{
		if invitation.Code == inviteCode{
			invitationLink = &invitation
			break
		}
	}

	if invitationLink == nil{
		return "",fmt.Errorf("linked invite not found",)
	}

	if invitationLink.ApprovedCount >= invitationLink.MaxApprovedJoins{
		return "", fmt.Errorf("invite limit reached (%d/%d approved joins). Create a new invite",invitationLink.ApprovedCount, invitationLink.MaxApprovedJoins)
	}

	// if time.Now().After(invitationLink.ExpiresAt){
	// 	return "",fmt.Errorf("cannot approve: the invite link used has expired")
	// }

	filter=bson.M{
		"slug":             slug,
		"joinRequests._id": requestID,
		"invites.code":     inviteCode,
	}

	update:=bson.M{
		"$addToSet": bson.M{"members": member},
		"$inc":      bson.M{"invites.$.approvedCount": 1},
		"$set": bson.M{
			"joinRequests.$[req].status":     "approved",
			"joinRequests.$[req].reviewedAt": time.Now(),
			"joinRequests.$[req].reviewedBy": reviewerId,
			"updatedAt":                      time.Now(),
		},
	}

	opts := options.UpdateOne().SetArrayFilters([]interface{}{bson.M{"req._id": requestID}})

	result,err:=repo.db.UpdateOne(ctx,filter,update,opts)
	if err != nil{
		return "",fmt.Errorf("error approving join request: %w",err)
	}
	if result.MatchedCount == 0 {
		return "",fmt.Errorf("request or workspace not found")
	}
	return "Approved the user to the workspace",nil
}

func (repo *WorkspaceRepositoryImpl)RejectRequest(ctx context.Context,slug string,inviteCode string,requestID bson.ObjectID,reviewerId string)(string,error){


	filter:=bson.M{
		"slug":             slug,
		"joinRequests._id": requestID,
	}

	update:=bson.M{
		"$set": bson.M{
			"joinRequests.$[req].status":     "rejected",
			"joinRequests.$[req].reviewedAt": time.Now(),
			"joinRequests.$[req].reviewedBy": reviewerId,
			"updatedAt":                      time.Now(),
		},
	}

	opts := options.UpdateOne().SetArrayFilters([]interface{}{bson.M{"req._id": requestID}})

	result,err:=repo.db.UpdateOne(ctx,filter,update,opts)
	if err != nil{
		return "",fmt.Errorf("error rejecting join request: %w",err)
	}
	if result.MatchedCount == 0 {
		return "",fmt.Errorf("request or workspace not found")
	}
	return "Rejected the user to the workspace",nil
}

func (repo *WorkspaceRepositoryImpl)UpdateWorkspace(ctx context.Context,slug string, mData *schemas.WorkSpace) (*schemas.WorkSpace,error){
	
	filter:=bson.M{
		"slug":slug,
	}

	setFields := bson.M{"updatedAt":time.Now(),}
	if mData.Name != ""{
		setFields["name"] = mData.Name
	}
	if mData.Description != ""{
		setFields["description"] = mData.Description
	}

	update:=bson.M{
		"$set":setFields,
	}

	options:=options.FindOneAndUpdate().SetReturnDocument(options.After) // it is used so that it will return the updated workspace back

	var mWorkspace schemas.WorkSpace

	err:=repo.db.FindOneAndUpdate(ctx,filter,update,options).Decode(&mWorkspace)
	if err!=nil{
		if err == mongo.ErrNoDocuments{
			return nil ,fmt.Errorf("workspace not found")
		}
		return nil , fmt.Errorf("Error occured in updating the workspace %w",err)
	}
	
	return &mWorkspace ,nil
}

func(repo *WorkspaceRepositoryImpl) RemoveMemberFromWorkspace(ctx context.Context,slug string,memberId string)(bool ,error){

	filter:=bson.M{
		"slug":slug,
	}

	update:=bson.M{
		"$pull":bson.M{

			"members":bson.M{
				"userId":memberId,
			},
		},
		"$set": bson.M{
        	"updatedAt": time.Now(),
    	},
	}
	result,err := repo.db.UpdateOne(ctx,filter,update)
	if err != nil{
		return false,fmt.Errorf("error removing  member: %w",err)
	}

	if result.MatchedCount == 0{
		return false,fmt.Errorf("workspace not found")
	}
	if result.ModifiedCount == 0 {
		return false, fmt.Errorf("user is not a member of this workspace")
	}

	return true,nil

}

func(repo *WorkspaceRepositoryImpl)DeleteWorkspace(ctx context.Context,slug string )(bool,error){
	
	filter:=bson.M{
		"slug":slug,
	}

	result,err := repo.db.DeleteOne(ctx,filter)
	if err!=nil{
		return false,fmt.Errorf("Error occured while deleting the workspaace %w",err)
	}
	if result.DeletedCount == 0{
		return false,fmt.Errorf("No workspae is present with this id to delete")
	}

	return true,nil
}