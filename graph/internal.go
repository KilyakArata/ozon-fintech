package graph

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"
	"sync"
	"sync/atomic"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/introspection"
	"github.com/KilyakArata/gqlgen-todos/graph/model"
	"github.com/KilyakArata/gqlgen-todos/models"
	gqlparser "github.com/vektah/gqlparser/v2"
	"github.com/vektah/gqlparser/v2/ast"
)

// NewExecutableSchema creates an ExecutableSchema from the ResolverRoot interface.
func NewExecutableSchema(cfg Config) graphql.ExecutableSchema {
	return &executableSchema{
		schema:     cfg.Schema,
		resolvers:  cfg.Resolvers,
		directives: cfg.Directives,
		complexity: cfg.Complexity,
	}
}

type Config struct {
	Schema     *ast.Schema
	Resolvers  ResolverRoot
	Directives DirectiveRoot
	Complexity ComplexityRoot
}

type ResolverRoot interface {
	Comment() CommentResolver
	Mutation() MutationResolver
	Post() PostResolver
	Query() QueryResolver
	Subscription() SubscriptionResolver
}

type DirectiveRoot struct {
}

type ComplexityRoot struct {
	Comment struct {
		Author   func(childComplexity int) int
		Body     func(childComplexity int) int
		Comments func(childComplexity int) int
		ID       func(childComplexity int) int
	}

	Mutation struct {
		CreateComment func(childComplexity int, input model.NewComment) int
		CreatePost    func(childComplexity int, input model.NewPost) int
		DeleteComment func(childComplexity int, id string) int
		DeletePost    func(childComplexity int, id string) int
		UpdateComment func(childComplexity int, id string, input model.UpdateComment) int
		UpdatePost    func(childComplexity int, id string, input model.UpdatePost) int
	}

	Post struct {
		Author   func(childComplexity int) int
		Body     func(childComplexity int) int
		Comments func(childComplexity int) int
		ID       func(childComplexity int) int
		Open     func(childComplexity int) int
	}

	Query struct {
		Comments func(childComplexity int, filter *model.CommentFilter, limit *int, offset *int) int
		Post     func(childComplexity int, id string) int
		Posts    func(childComplexity int) int
	}

	Subscription struct {
		CommentAdded func(childComplexity int, postID string) int
	}
}

type CommentResolver interface {
	Comments(ctx context.Context, obj *models.Comment) ([]*models.Comment, error)
}
type MutationResolver interface {
	CreateComment(ctx context.Context, input model.NewComment) (*models.Comment, error)
	UpdateComment(ctx context.Context, id string, input model.UpdateComment) (*models.Comment, error)
	DeleteComment(ctx context.Context, id string) (bool, error)
	CreatePost(ctx context.Context, input model.NewPost) (*models.Post, error)
	UpdatePost(ctx context.Context, id string, input model.UpdatePost) (*models.Post, error)
	DeletePost(ctx context.Context, id string) (bool, error)
}
type PostResolver interface {
	Comments(ctx context.Context, obj *models.Post) ([]*models.Comment, error)
}
type QueryResolver interface {
	Comments(ctx context.Context, filter *model.CommentFilter, limit *int, offset *int) ([]*models.Comment, error)
	Posts(ctx context.Context) ([]*models.Post, error)
	Post(ctx context.Context, id string) (*models.Post, error)
}
type SubscriptionResolver interface {
	CommentAdded(ctx context.Context, postID string) (<-chan *models.Comment, error)
}

type executableSchema struct {
	schema     *ast.Schema
	resolvers  ResolverRoot
	directives DirectiveRoot
	complexity ComplexityRoot
}

func (e *executableSchema) Schema() *ast.Schema {
	if e.schema != nil {
		return e.schema
	}
	return parsedSchema
}

func (e *executableSchema) Complexity(typeName, field string, childComplexity int, rawArgs map[string]interface{}) (int, bool) {
	ec := executionContext{nil, e, 0, 0, nil}
	_ = ec
	switch typeName + "." + field {

	case "Comment.author":
		if e.complexity.Comment.Author == nil {
			break
		}

		return e.complexity.Comment.Author(childComplexity), true

	case "Comment.body":
		if e.complexity.Comment.Body == nil {
			break
		}

		return e.complexity.Comment.Body(childComplexity), true

	case "Comment.comments":
		if e.complexity.Comment.Comments == nil {
			break
		}

		return e.complexity.Comment.Comments(childComplexity), true

	case "Comment.id":
		if e.complexity.Comment.ID == nil {
			break
		}

		return e.complexity.Comment.ID(childComplexity), true

	case "Mutation.createComment":
		if e.complexity.Mutation.CreateComment == nil {
			break
		}

		args, err := ec.field_Mutation_createComment_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.CreateComment(childComplexity, args["input"].(model.NewComment)), true

	case "Mutation.createPost":
		if e.complexity.Mutation.CreatePost == nil {
			break
		}

		args, err := ec.field_Mutation_createPost_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.CreatePost(childComplexity, args["input"].(model.NewPost)), true

	case "Mutation.deleteComment":
		if e.complexity.Mutation.DeleteComment == nil {
			break
		}

		args, err := ec.field_Mutation_deleteComment_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.DeleteComment(childComplexity, args["id"].(string)), true

	case "Mutation.deletePost":
		if e.complexity.Mutation.DeletePost == nil {
			break
		}

		args, err := ec.field_Mutation_deletePost_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.DeletePost(childComplexity, args["id"].(string)), true

	case "Mutation.updateComment":
		if e.complexity.Mutation.UpdateComment == nil {
			break
		}

		args, err := ec.field_Mutation_updateComment_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.UpdateComment(childComplexity, args["id"].(string), args["input"].(model.UpdateComment)), true

	case "Mutation.updatePost":
		if e.complexity.Mutation.UpdatePost == nil {
			break
		}

		args, err := ec.field_Mutation_updatePost_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Mutation.UpdatePost(childComplexity, args["id"].(string), args["input"].(model.UpdatePost)), true

	case "Post.author":
		if e.complexity.Post.Author == nil {
			break
		}

		return e.complexity.Post.Author(childComplexity), true

	case "Post.body":
		if e.complexity.Post.Body == nil {
			break
		}

		return e.complexity.Post.Body(childComplexity), true

	case "Post.comments":
		if e.complexity.Post.Comments == nil {
			break
		}

		return e.complexity.Post.Comments(childComplexity), true

	case "Post.id":
		if e.complexity.Post.ID == nil {
			break
		}

		return e.complexity.Post.ID(childComplexity), true

	case "Post.open":
		if e.complexity.Post.Open == nil {
			break
		}

		return e.complexity.Post.Open(childComplexity), true

	case "Query.comments":
		if e.complexity.Query.Comments == nil {
			break
		}

		args, err := ec.field_Query_comments_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Query.Comments(childComplexity, args["filter"].(*model.CommentFilter), args["limit"].(*int), args["offset"].(*int)), true

	case "Query.post":
		if e.complexity.Query.Post == nil {
			break
		}

		args, err := ec.field_Query_post_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Query.Post(childComplexity, args["id"].(string)), true

	case "Query.posts":
		if e.complexity.Query.Posts == nil {
			break
		}

		return e.complexity.Query.Posts(childComplexity), true

	case "Subscription.commentAdded":
		if e.complexity.Subscription.CommentAdded == nil {
			break
		}

		args, err := ec.field_Subscription_commentAdded_args(context.TODO(), rawArgs)
		if err != nil {
			return 0, false
		}

		return e.complexity.Subscription.CommentAdded(childComplexity, args["postId"].(string)), true

	}
	return 0, false
}

func (e *executableSchema) Exec(ctx context.Context) graphql.ResponseHandler {
	rc := graphql.GetOperationContext(ctx)
	ec := executionContext{rc, e, 0, 0, make(chan graphql.DeferredResult)}
	inputUnmarshalMap := graphql.BuildUnmarshalerMap(
		ec.unmarshalInputCommentFilter,
		ec.unmarshalInputNewComment,
		ec.unmarshalInputNewPost,
		ec.unmarshalInputUpdateComment,
		ec.unmarshalInputUpdatePost,
	)
	first := true

	switch rc.Operation.Operation {
	case ast.Query:
		return func(ctx context.Context) *graphql.Response {
			var response graphql.Response
			var data graphql.Marshaler
			if first {
				first = false
				ctx = graphql.WithUnmarshalerMap(ctx, inputUnmarshalMap)
				data = ec._Query(ctx, rc.Operation.SelectionSet)
			} else {
				if atomic.LoadInt32(&ec.pendingDeferred) > 0 {
					result := <-ec.deferredResults
					atomic.AddInt32(&ec.pendingDeferred, -1)
					data = result.Result
					response.Path = result.Path
					response.Label = result.Label
					response.Errors = result.Errors
				} else {
					return nil
				}
			}
			var buf bytes.Buffer
			data.MarshalGQL(&buf)
			response.Data = buf.Bytes()
			if atomic.LoadInt32(&ec.deferred) > 0 {
				hasNext := atomic.LoadInt32(&ec.pendingDeferred) > 0
				response.HasNext = &hasNext
			}

			return &response
		}
	case ast.Mutation:
		return func(ctx context.Context) *graphql.Response {
			if !first {
				return nil
			}
			first = false
			ctx = graphql.WithUnmarshalerMap(ctx, inputUnmarshalMap)
			data := ec._Mutation(ctx, rc.Operation.SelectionSet)
			var buf bytes.Buffer
			data.MarshalGQL(&buf)

			return &graphql.Response{
				Data: buf.Bytes(),
			}
		}
	case ast.Subscription:
		next := ec._Subscription(ctx, rc.Operation.SelectionSet)

		var buf bytes.Buffer
		return func(ctx context.Context) *graphql.Response {
			buf.Reset()
			data := next(ctx)

			if data == nil {
				return nil
			}
			data.MarshalGQL(&buf)

			return &graphql.Response{
				Data: buf.Bytes(),
			}
		}

	default:
		return graphql.OneShot(graphql.ErrorResponse(ctx, "unsupported GraphQL operation"))
	}
}

type executionContext struct {
	*graphql.OperationContext
	*executableSchema
	deferred        int32
	pendingDeferred int32
	deferredResults chan graphql.DeferredResult
}

func (ec *executionContext) processDeferredGroup(dg graphql.DeferredGroup) {
	atomic.AddInt32(&ec.pendingDeferred, 1)
	go func() {
		ctx := graphql.WithFreshResponseContext(dg.Context)
		dg.FieldSet.Dispatch(ctx)
		ds := graphql.DeferredResult{
			Path:   dg.Path,
			Label:  dg.Label,
			Result: dg.FieldSet,
			Errors: graphql.GetErrors(ctx),
		}
		// null fields should bubble up
		if dg.FieldSet.Invalids > 0 {
			ds.Result = graphql.Null
		}
		ec.deferredResults <- ds
	}()
}

func (ec *executionContext) introspectSchema() (*introspection.Schema, error) {
	if ec.DisableIntrospection {
		return nil, errors.New("introspection disabled")
	}
	return introspection.WrapSchema(ec.Schema()), nil
}

func (ec *executionContext) introspectType(name string) (*introspection.Type, error) {
	if ec.DisableIntrospection {
		return nil, errors.New("introspection disabled")
	}
	return introspection.WrapTypeFromDef(ec.Schema(), ec.Schema().Types[name]), nil
}

var sources = []*ast.Source{
	{Name: "../schema.graphqls", Input: `type Post{
    id: ID!
    author: String!
    body: String!
    open: Boolean!

    comments: [Comment!]
}

type Comment{
    id: ID!
    author: String!
    body: String!

    comments: [Comment!]
}

input NewComment{
    author: String!
    body: String!
    parentPost: ID!
    parentComment: ID
}

input NewPost{
    author: String!
    body: String!
    open: Boolean!
}

input CommentFilter {
    author: String
    body: String
}

input UpdatePost {
    body: String
    open: Boolean
}

input UpdateComment {
    body: String
}

type Query{
    comments(filter: CommentFilter, limit: Int = 10, offset: Int = 0): [Comment!]!
    posts: [Post!]!
    post(id:  ID!): Post!
}

type Mutation {
    createComment(input: NewComment!): Comment!
    updateComment(id: ID!, input: UpdateComment!): Comment!
    deleteComment(id: ID!): Boolean!
    createPost(input: NewPost!): Post!
    updatePost(id: ID!, input: UpdatePost!): Post!
    deletePost(id: ID!): Boolean!
}

type Subscription {
    commentAdded(postId: ID!): Comment!
}`, BuiltIn: false},
}
var parsedSchema = gqlparser.MustLoadSchema(sources...)

func (ec *executionContext) field_Mutation_createComment_args(ctx context.Context, rawArgs map[string]interface{}) (map[string]interface{}, error) {
	var err error
	args := map[string]interface{}{}
	var arg0 model.NewComment
	if tmp, ok := rawArgs["input"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("input"))
		arg0, err = ec.unmarshalNNewComment2githubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋgraphᚋmodelᚐNewComment(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["input"] = arg0
	return args, nil
}

func (ec *executionContext) field_Mutation_createPost_args(ctx context.Context, rawArgs map[string]interface{}) (map[string]interface{}, error) {
	var err error
	args := map[string]interface{}{}
	var arg0 model.NewPost
	if tmp, ok := rawArgs["input"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("input"))
		arg0, err = ec.unmarshalNNewPost2githubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋgraphᚋmodelᚐNewPost(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["input"] = arg0
	return args, nil
}

func (ec *executionContext) field_Mutation_deleteComment_args(ctx context.Context, rawArgs map[string]interface{}) (map[string]interface{}, error) {
	var err error
	args := map[string]interface{}{}
	var arg0 string
	if tmp, ok := rawArgs["id"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("id"))
		arg0, err = ec.unmarshalNID2string(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["id"] = arg0
	return args, nil
}

func (ec *executionContext) field_Mutation_deletePost_args(ctx context.Context, rawArgs map[string]interface{}) (map[string]interface{}, error) {
	var err error
	args := map[string]interface{}{}
	var arg0 string
	if tmp, ok := rawArgs["id"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("id"))
		arg0, err = ec.unmarshalNID2string(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["id"] = arg0
	return args, nil
}

func (ec *executionContext) field_Mutation_updateComment_args(ctx context.Context, rawArgs map[string]interface{}) (map[string]interface{}, error) {
	var err error
	args := map[string]interface{}{}
	var arg0 string
	if tmp, ok := rawArgs["id"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("id"))
		arg0, err = ec.unmarshalNID2string(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["id"] = arg0
	var arg1 model.UpdateComment
	if tmp, ok := rawArgs["input"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("input"))
		arg1, err = ec.unmarshalNUpdateComment2githubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋgraphᚋmodelᚐUpdateComment(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["input"] = arg1
	return args, nil
}

func (ec *executionContext) field_Mutation_updatePost_args(ctx context.Context, rawArgs map[string]interface{}) (map[string]interface{}, error) {
	var err error
	args := map[string]interface{}{}
	var arg0 string
	if tmp, ok := rawArgs["id"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("id"))
		arg0, err = ec.unmarshalNID2string(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["id"] = arg0
	var arg1 model.UpdatePost
	if tmp, ok := rawArgs["input"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("input"))
		arg1, err = ec.unmarshalNUpdatePost2githubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋgraphᚋmodelᚐUpdatePost(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["input"] = arg1
	return args, nil
}

func (ec *executionContext) field_Query___type_args(ctx context.Context, rawArgs map[string]interface{}) (map[string]interface{}, error) {
	var err error
	args := map[string]interface{}{}
	var arg0 string
	if tmp, ok := rawArgs["name"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("name"))
		arg0, err = ec.unmarshalNString2string(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["name"] = arg0
	return args, nil
}

func (ec *executionContext) field_Query_comments_args(ctx context.Context, rawArgs map[string]interface{}) (map[string]interface{}, error) {
	var err error
	args := map[string]interface{}{}
	var arg0 *model.CommentFilter
	if tmp, ok := rawArgs["filter"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("filter"))
		arg0, err = ec.unmarshalOCommentFilter2ᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋgraphᚋmodelᚐCommentFilter(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["filter"] = arg0
	var arg1 *int
	if tmp, ok := rawArgs["limit"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("limit"))
		arg1, err = ec.unmarshalOInt2ᚖint(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["limit"] = arg1
	var arg2 *int
	if tmp, ok := rawArgs["offset"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("offset"))
		arg2, err = ec.unmarshalOInt2ᚖint(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["offset"] = arg2
	return args, nil
}

func (ec *executionContext) field_Query_post_args(ctx context.Context, rawArgs map[string]interface{}) (map[string]interface{}, error) {
	var err error
	args := map[string]interface{}{}
	var arg0 string
	if tmp, ok := rawArgs["id"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("id"))
		arg0, err = ec.unmarshalNID2string(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["id"] = arg0
	return args, nil
}

func (ec *executionContext) field_Subscription_commentAdded_args(ctx context.Context, rawArgs map[string]interface{}) (map[string]interface{}, error) {
	var err error
	args := map[string]interface{}{}
	var arg0 string
	if tmp, ok := rawArgs["postId"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("postId"))
		arg0, err = ec.unmarshalNID2string(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["postId"] = arg0
	return args, nil
}

func (ec *executionContext) field___Type_enumValues_args(ctx context.Context, rawArgs map[string]interface{}) (map[string]interface{}, error) {
	var err error
	args := map[string]interface{}{}
	var arg0 bool
	if tmp, ok := rawArgs["includeDeprecated"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("includeDeprecated"))
		arg0, err = ec.unmarshalOBoolean2bool(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["includeDeprecated"] = arg0
	return args, nil
}

func (ec *executionContext) field___Type_fields_args(ctx context.Context, rawArgs map[string]interface{}) (map[string]interface{}, error) {
	var err error
	args := map[string]interface{}{}
	var arg0 bool
	if tmp, ok := rawArgs["includeDeprecated"]; ok {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("includeDeprecated"))
		arg0, err = ec.unmarshalOBoolean2bool(ctx, tmp)
		if err != nil {
			return nil, err
		}
	}
	args["includeDeprecated"] = arg0
	return args, nil
}

func (ec *executionContext) _Comment_id(ctx context.Context, field graphql.CollectedField, obj *models.Comment) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Comment_id(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.ID, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(string)
	fc.Result = res
	return ec.marshalNID2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Comment_id(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Comment",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type ID does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Comment_author(ctx context.Context, field graphql.CollectedField, obj *models.Comment) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Comment_author(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Author, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(string)
	fc.Result = res
	return ec.marshalNString2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Comment_author(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Comment",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Comment_body(ctx context.Context, field graphql.CollectedField, obj *models.Comment) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Comment_body(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Body, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(string)
	fc.Result = res
	return ec.marshalNString2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Comment_body(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Comment",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Comment_comments(ctx context.Context, field graphql.CollectedField, obj *models.Comment) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Comment_comments(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Comment().Comments(rctx, obj)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.([]*models.Comment)
	fc.Result = res
	return ec.marshalOComment2ᚕᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐCommentᚄ(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Comment_comments(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Comment",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "id":
				return ec.fieldContext_Comment_id(ctx, field)
			case "author":
				return ec.fieldContext_Comment_author(ctx, field)
			case "body":
				return ec.fieldContext_Comment_body(ctx, field)
			case "comments":
				return ec.fieldContext_Comment_comments(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type Comment", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) _Mutation_createComment(ctx context.Context, field graphql.CollectedField) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Mutation_createComment(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Mutation().CreateComment(rctx, fc.Args["input"].(model.NewComment))
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*models.Comment)
	fc.Result = res
	return ec.marshalNComment2ᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐComment(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Mutation_createComment(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Mutation",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "id":
				return ec.fieldContext_Comment_id(ctx, field)
			case "author":
				return ec.fieldContext_Comment_author(ctx, field)
			case "body":
				return ec.fieldContext_Comment_body(ctx, field)
			case "comments":
				return ec.fieldContext_Comment_comments(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type Comment", field.Name)
		},
	}
	defer func() {
		if r := recover(); r != nil {
			err = ec.Recover(ctx, r)
			ec.Error(ctx, err)
		}
	}()
	ctx = graphql.WithFieldContext(ctx, fc)
	if fc.Args, err = ec.field_Mutation_createComment_args(ctx, field.ArgumentMap(ec.Variables)); err != nil {
		ec.Error(ctx, err)
		return fc, err
	}
	return fc, nil
}

func (ec *executionContext) _Mutation_updateComment(ctx context.Context, field graphql.CollectedField) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Mutation_updateComment(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Mutation().UpdateComment(rctx, fc.Args["id"].(string), fc.Args["input"].(model.UpdateComment))
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*models.Comment)
	fc.Result = res
	return ec.marshalNComment2ᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐComment(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Mutation_updateComment(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Mutation",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "id":
				return ec.fieldContext_Comment_id(ctx, field)
			case "author":
				return ec.fieldContext_Comment_author(ctx, field)
			case "body":
				return ec.fieldContext_Comment_body(ctx, field)
			case "comments":
				return ec.fieldContext_Comment_comments(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type Comment", field.Name)
		},
	}
	defer func() {
		if r := recover(); r != nil {
			err = ec.Recover(ctx, r)
			ec.Error(ctx, err)
		}
	}()
	ctx = graphql.WithFieldContext(ctx, fc)
	if fc.Args, err = ec.field_Mutation_updateComment_args(ctx, field.ArgumentMap(ec.Variables)); err != nil {
		ec.Error(ctx, err)
		return fc, err
	}
	return fc, nil
}

func (ec *executionContext) _Mutation_deleteComment(ctx context.Context, field graphql.CollectedField) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Mutation_deleteComment(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Mutation().DeleteComment(rctx, fc.Args["id"].(string))
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(bool)
	fc.Result = res
	return ec.marshalNBoolean2bool(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Mutation_deleteComment(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Mutation",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type Boolean does not have child fields")
		},
	}
	defer func() {
		if r := recover(); r != nil {
			err = ec.Recover(ctx, r)
			ec.Error(ctx, err)
		}
	}()
	ctx = graphql.WithFieldContext(ctx, fc)
	if fc.Args, err = ec.field_Mutation_deleteComment_args(ctx, field.ArgumentMap(ec.Variables)); err != nil {
		ec.Error(ctx, err)
		return fc, err
	}
	return fc, nil
}

func (ec *executionContext) _Mutation_createPost(ctx context.Context, field graphql.CollectedField) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Mutation_createPost(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Mutation().CreatePost(rctx, fc.Args["input"].(model.NewPost))
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*models.Post)
	fc.Result = res
	return ec.marshalNPost2ᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐPost(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Mutation_createPost(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Mutation",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "id":
				return ec.fieldContext_Post_id(ctx, field)
			case "author":
				return ec.fieldContext_Post_author(ctx, field)
			case "body":
				return ec.fieldContext_Post_body(ctx, field)
			case "open":
				return ec.fieldContext_Post_open(ctx, field)
			case "comments":
				return ec.fieldContext_Post_comments(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type Post", field.Name)
		},
	}
	defer func() {
		if r := recover(); r != nil {
			err = ec.Recover(ctx, r)
			ec.Error(ctx, err)
		}
	}()
	ctx = graphql.WithFieldContext(ctx, fc)
	if fc.Args, err = ec.field_Mutation_createPost_args(ctx, field.ArgumentMap(ec.Variables)); err != nil {
		ec.Error(ctx, err)
		return fc, err
	}
	return fc, nil
}

func (ec *executionContext) _Mutation_updatePost(ctx context.Context, field graphql.CollectedField) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Mutation_updatePost(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Mutation().UpdatePost(rctx, fc.Args["id"].(string), fc.Args["input"].(model.UpdatePost))
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*models.Post)
	fc.Result = res
	return ec.marshalNPost2ᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐPost(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Mutation_updatePost(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Mutation",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "id":
				return ec.fieldContext_Post_id(ctx, field)
			case "author":
				return ec.fieldContext_Post_author(ctx, field)
			case "body":
				return ec.fieldContext_Post_body(ctx, field)
			case "open":
				return ec.fieldContext_Post_open(ctx, field)
			case "comments":
				return ec.fieldContext_Post_comments(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type Post", field.Name)
		},
	}
	defer func() {
		if r := recover(); r != nil {
			err = ec.Recover(ctx, r)
			ec.Error(ctx, err)
		}
	}()
	ctx = graphql.WithFieldContext(ctx, fc)
	if fc.Args, err = ec.field_Mutation_updatePost_args(ctx, field.ArgumentMap(ec.Variables)); err != nil {
		ec.Error(ctx, err)
		return fc, err
	}
	return fc, nil
}

func (ec *executionContext) _Mutation_deletePost(ctx context.Context, field graphql.CollectedField) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Mutation_deletePost(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Mutation().DeletePost(rctx, fc.Args["id"].(string))
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(bool)
	fc.Result = res
	return ec.marshalNBoolean2bool(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Mutation_deletePost(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Mutation",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type Boolean does not have child fields")
		},
	}
	defer func() {
		if r := recover(); r != nil {
			err = ec.Recover(ctx, r)
			ec.Error(ctx, err)
		}
	}()
	ctx = graphql.WithFieldContext(ctx, fc)
	if fc.Args, err = ec.field_Mutation_deletePost_args(ctx, field.ArgumentMap(ec.Variables)); err != nil {
		ec.Error(ctx, err)
		return fc, err
	}
	return fc, nil
}

func (ec *executionContext) _Post_id(ctx context.Context, field graphql.CollectedField, obj *models.Post) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Post_id(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.ID, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(string)
	fc.Result = res
	return ec.marshalNID2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Post_id(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Post",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type ID does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Post_author(ctx context.Context, field graphql.CollectedField, obj *models.Post) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Post_author(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Author, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(string)
	fc.Result = res
	return ec.marshalNString2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Post_author(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Post",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Post_body(ctx context.Context, field graphql.CollectedField, obj *models.Post) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Post_body(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Body, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(string)
	fc.Result = res
	return ec.marshalNString2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Post_body(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Post",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Post_open(ctx context.Context, field graphql.CollectedField, obj *models.Post) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Post_open(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Open, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(bool)
	fc.Result = res
	return ec.marshalNBoolean2bool(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Post_open(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Post",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type Boolean does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) _Post_comments(ctx context.Context, field graphql.CollectedField, obj *models.Post) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Post_comments(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Post().Comments(rctx, obj)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.([]*models.Comment)
	fc.Result = res
	return ec.marshalOComment2ᚕᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐCommentᚄ(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Post_comments(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Post",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "id":
				return ec.fieldContext_Comment_id(ctx, field)
			case "author":
				return ec.fieldContext_Comment_author(ctx, field)
			case "body":
				return ec.fieldContext_Comment_body(ctx, field)
			case "comments":
				return ec.fieldContext_Comment_comments(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type Comment", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) _Query_comments(ctx context.Context, field graphql.CollectedField) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Query_comments(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Query().Comments(rctx, fc.Args["filter"].(*model.CommentFilter), fc.Args["limit"].(*int), fc.Args["offset"].(*int))
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.([]*models.Comment)
	fc.Result = res
	return ec.marshalNComment2ᚕᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐCommentᚄ(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Query_comments(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Query",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "id":
				return ec.fieldContext_Comment_id(ctx, field)
			case "author":
				return ec.fieldContext_Comment_author(ctx, field)
			case "body":
				return ec.fieldContext_Comment_body(ctx, field)
			case "comments":
				return ec.fieldContext_Comment_comments(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type Comment", field.Name)
		},
	}
	defer func() {
		if r := recover(); r != nil {
			err = ec.Recover(ctx, r)
			ec.Error(ctx, err)
		}
	}()
	ctx = graphql.WithFieldContext(ctx, fc)
	if fc.Args, err = ec.field_Query_comments_args(ctx, field.ArgumentMap(ec.Variables)); err != nil {
		ec.Error(ctx, err)
		return fc, err
	}
	return fc, nil
}

func (ec *executionContext) _Query_posts(ctx context.Context, field graphql.CollectedField) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Query_posts(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Query().Posts(rctx)
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.([]*models.Post)
	fc.Result = res
	return ec.marshalNPost2ᚕᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐPostᚄ(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Query_posts(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Query",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "id":
				return ec.fieldContext_Post_id(ctx, field)
			case "author":
				return ec.fieldContext_Post_author(ctx, field)
			case "body":
				return ec.fieldContext_Post_body(ctx, field)
			case "open":
				return ec.fieldContext_Post_open(ctx, field)
			case "comments":
				return ec.fieldContext_Post_comments(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type Post", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) _Query_post(ctx context.Context, field graphql.CollectedField) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Query_post(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Query().Post(rctx, fc.Args["id"].(string))
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*models.Post)
	fc.Result = res
	return ec.marshalNPost2ᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐPost(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Query_post(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Query",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "id":
				return ec.fieldContext_Post_id(ctx, field)
			case "author":
				return ec.fieldContext_Post_author(ctx, field)
			case "body":
				return ec.fieldContext_Post_body(ctx, field)
			case "open":
				return ec.fieldContext_Post_open(ctx, field)
			case "comments":
				return ec.fieldContext_Post_comments(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type Post", field.Name)
		},
	}
	defer func() {
		if r := recover(); r != nil {
			err = ec.Recover(ctx, r)
			ec.Error(ctx, err)
		}
	}()
	ctx = graphql.WithFieldContext(ctx, fc)
	if fc.Args, err = ec.field_Query_post_args(ctx, field.ArgumentMap(ec.Variables)); err != nil {
		ec.Error(ctx, err)
		return fc, err
	}
	return fc, nil
}

func (ec *executionContext) _Query___type(ctx context.Context, field graphql.CollectedField) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Query___type(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.introspectType(fc.Args["name"].(string))
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.(*introspection.Type)
	fc.Result = res
	return ec.marshalO__Type2ᚖgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐType(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Query___type(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Query",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "kind":
				return ec.fieldContext___Type_kind(ctx, field)
			case "name":
				return ec.fieldContext___Type_name(ctx, field)
			case "description":
				return ec.fieldContext___Type_description(ctx, field)
			case "fields":
				return ec.fieldContext___Type_fields(ctx, field)
			case "interfaces":
				return ec.fieldContext___Type_interfaces(ctx, field)
			case "possibleTypes":
				return ec.fieldContext___Type_possibleTypes(ctx, field)
			case "enumValues":
				return ec.fieldContext___Type_enumValues(ctx, field)
			case "inputFields":
				return ec.fieldContext___Type_inputFields(ctx, field)
			case "ofType":
				return ec.fieldContext___Type_ofType(ctx, field)
			case "specifiedByURL":
				return ec.fieldContext___Type_specifiedByURL(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type __Type", field.Name)
		},
	}
	defer func() {
		if r := recover(); r != nil {
			err = ec.Recover(ctx, r)
			ec.Error(ctx, err)
		}
	}()
	ctx = graphql.WithFieldContext(ctx, fc)
	if fc.Args, err = ec.field_Query___type_args(ctx, field.ArgumentMap(ec.Variables)); err != nil {
		ec.Error(ctx, err)
		return fc, err
	}
	return fc, nil
}

func (ec *executionContext) _Query___schema(ctx context.Context, field graphql.CollectedField) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext_Query___schema(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.introspectSchema()
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.(*introspection.Schema)
	fc.Result = res
	return ec.marshalO__Schema2ᚖgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐSchema(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext_Query___schema(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Query",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "description":
				return ec.fieldContext___Schema_description(ctx, field)
			case "types":
				return ec.fieldContext___Schema_types(ctx, field)
			case "queryType":
				return ec.fieldContext___Schema_queryType(ctx, field)
			case "mutationType":
				return ec.fieldContext___Schema_mutationType(ctx, field)
			case "subscriptionType":
				return ec.fieldContext___Schema_subscriptionType(ctx, field)
			case "directives":
				return ec.fieldContext___Schema_directives(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type __Schema", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) _Subscription_commentAdded(ctx context.Context, field graphql.CollectedField) (ret func(ctx context.Context) graphql.Marshaler) {
	fc, err := ec.fieldContext_Subscription_commentAdded(ctx, field)
	if err != nil {
		return nil
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = nil
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return ec.resolvers.Subscription().CommentAdded(rctx, fc.Args["postId"].(string))
	})
	if err != nil {
		ec.Error(ctx, err)
		return nil
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return nil
	}
	return func(ctx context.Context) graphql.Marshaler {
		select {
		case res, ok := <-resTmp.(<-chan *models.Comment):
			if !ok {
				return nil
			}
			return graphql.WriterFunc(func(w io.Writer) {
				w.Write([]byte{'{'})
				graphql.MarshalString(field.Alias).MarshalGQL(w)
				w.Write([]byte{':'})
				ec.marshalNComment2ᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐComment(ctx, field.Selections, res).MarshalGQL(w)
				w.Write([]byte{'}'})
			})
		case <-ctx.Done():
			return nil
		}
	}
}

func (ec *executionContext) fieldContext_Subscription_commentAdded(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "Subscription",
		Field:      field,
		IsMethod:   true,
		IsResolver: true,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "id":
				return ec.fieldContext_Comment_id(ctx, field)
			case "author":
				return ec.fieldContext_Comment_author(ctx, field)
			case "body":
				return ec.fieldContext_Comment_body(ctx, field)
			case "comments":
				return ec.fieldContext_Comment_comments(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type Comment", field.Name)
		},
	}
	defer func() {
		if r := recover(); r != nil {
			err = ec.Recover(ctx, r)
			ec.Error(ctx, err)
		}
	}()
	ctx = graphql.WithFieldContext(ctx, fc)
	if fc.Args, err = ec.field_Subscription_commentAdded_args(ctx, field.ArgumentMap(ec.Variables)); err != nil {
		ec.Error(ctx, err)
		return fc, err
	}
	return fc, nil
}

func (ec *executionContext) ___Directive_name(ctx context.Context, field graphql.CollectedField, obj *introspection.Directive) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Directive_name(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Name, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(string)
	fc.Result = res
	return ec.marshalNString2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Directive_name(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Directive",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Directive_description(ctx context.Context, field graphql.CollectedField, obj *introspection.Directive) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Directive_description(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Description(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.(*string)
	fc.Result = res
	return ec.marshalOString2ᚖstring(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Directive_description(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Directive",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Directive_locations(ctx context.Context, field graphql.CollectedField, obj *introspection.Directive) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Directive_locations(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Locations, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.([]string)
	fc.Result = res
	return ec.marshalN__DirectiveLocation2ᚕstringᚄ(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Directive_locations(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Directive",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type __DirectiveLocation does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Directive_args(ctx context.Context, field graphql.CollectedField, obj *introspection.Directive) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Directive_args(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Args, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.([]introspection.InputValue)
	fc.Result = res
	return ec.marshalN__InputValue2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐInputValueᚄ(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Directive_args(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Directive",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "name":
				return ec.fieldContext___InputValue_name(ctx, field)
			case "description":
				return ec.fieldContext___InputValue_description(ctx, field)
			case "type":
				return ec.fieldContext___InputValue_type(ctx, field)
			case "defaultValue":
				return ec.fieldContext___InputValue_defaultValue(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type __InputValue", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Directive_isRepeatable(ctx context.Context, field graphql.CollectedField, obj *introspection.Directive) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Directive_isRepeatable(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.IsRepeatable, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(bool)
	fc.Result = res
	return ec.marshalNBoolean2bool(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Directive_isRepeatable(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Directive",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type Boolean does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) ___EnumValue_name(ctx context.Context, field graphql.CollectedField, obj *introspection.EnumValue) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___EnumValue_name(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Name, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(string)
	fc.Result = res
	return ec.marshalNString2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___EnumValue_name(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__EnumValue",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) ___EnumValue_description(ctx context.Context, field graphql.CollectedField, obj *introspection.EnumValue) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___EnumValue_description(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Description(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.(*string)
	fc.Result = res
	return ec.marshalOString2ᚖstring(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___EnumValue_description(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__EnumValue",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) ___EnumValue_isDeprecated(ctx context.Context, field graphql.CollectedField, obj *introspection.EnumValue) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___EnumValue_isDeprecated(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.IsDeprecated(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(bool)
	fc.Result = res
	return ec.marshalNBoolean2bool(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___EnumValue_isDeprecated(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__EnumValue",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type Boolean does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) ___EnumValue_deprecationReason(ctx context.Context, field graphql.CollectedField, obj *introspection.EnumValue) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___EnumValue_deprecationReason(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.DeprecationReason(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.(*string)
	fc.Result = res
	return ec.marshalOString2ᚖstring(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___EnumValue_deprecationReason(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__EnumValue",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Field_name(ctx context.Context, field graphql.CollectedField, obj *introspection.Field) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Field_name(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Name, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(string)
	fc.Result = res
	return ec.marshalNString2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Field_name(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Field",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Field_description(ctx context.Context, field graphql.CollectedField, obj *introspection.Field) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Field_description(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Description(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.(*string)
	fc.Result = res
	return ec.marshalOString2ᚖstring(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Field_description(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Field",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Field_args(ctx context.Context, field graphql.CollectedField, obj *introspection.Field) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Field_args(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Args, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.([]introspection.InputValue)
	fc.Result = res
	return ec.marshalN__InputValue2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐInputValueᚄ(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Field_args(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Field",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "name":
				return ec.fieldContext___InputValue_name(ctx, field)
			case "description":
				return ec.fieldContext___InputValue_description(ctx, field)
			case "type":
				return ec.fieldContext___InputValue_type(ctx, field)
			case "defaultValue":
				return ec.fieldContext___InputValue_defaultValue(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type __InputValue", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Field_type(ctx context.Context, field graphql.CollectedField, obj *introspection.Field) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Field_type(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Type, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*introspection.Type)
	fc.Result = res
	return ec.marshalN__Type2ᚖgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐType(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Field_type(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Field",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "kind":
				return ec.fieldContext___Type_kind(ctx, field)
			case "name":
				return ec.fieldContext___Type_name(ctx, field)
			case "description":
				return ec.fieldContext___Type_description(ctx, field)
			case "fields":
				return ec.fieldContext___Type_fields(ctx, field)
			case "interfaces":
				return ec.fieldContext___Type_interfaces(ctx, field)
			case "possibleTypes":
				return ec.fieldContext___Type_possibleTypes(ctx, field)
			case "enumValues":
				return ec.fieldContext___Type_enumValues(ctx, field)
			case "inputFields":
				return ec.fieldContext___Type_inputFields(ctx, field)
			case "ofType":
				return ec.fieldContext___Type_ofType(ctx, field)
			case "specifiedByURL":
				return ec.fieldContext___Type_specifiedByURL(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type __Type", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Field_isDeprecated(ctx context.Context, field graphql.CollectedField, obj *introspection.Field) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Field_isDeprecated(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.IsDeprecated(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(bool)
	fc.Result = res
	return ec.marshalNBoolean2bool(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Field_isDeprecated(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Field",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type Boolean does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Field_deprecationReason(ctx context.Context, field graphql.CollectedField, obj *introspection.Field) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Field_deprecationReason(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.DeprecationReason(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.(*string)
	fc.Result = res
	return ec.marshalOString2ᚖstring(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Field_deprecationReason(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Field",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) ___InputValue_name(ctx context.Context, field graphql.CollectedField, obj *introspection.InputValue) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___InputValue_name(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Name, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(string)
	fc.Result = res
	return ec.marshalNString2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___InputValue_name(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__InputValue",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) ___InputValue_description(ctx context.Context, field graphql.CollectedField, obj *introspection.InputValue) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___InputValue_description(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Description(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.(*string)
	fc.Result = res
	return ec.marshalOString2ᚖstring(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___InputValue_description(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__InputValue",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) ___InputValue_type(ctx context.Context, field graphql.CollectedField, obj *introspection.InputValue) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___InputValue_type(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Type, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*introspection.Type)
	fc.Result = res
	return ec.marshalN__Type2ᚖgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐType(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___InputValue_type(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__InputValue",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "kind":
				return ec.fieldContext___Type_kind(ctx, field)
			case "name":
				return ec.fieldContext___Type_name(ctx, field)
			case "description":
				return ec.fieldContext___Type_description(ctx, field)
			case "fields":
				return ec.fieldContext___Type_fields(ctx, field)
			case "interfaces":
				return ec.fieldContext___Type_interfaces(ctx, field)
			case "possibleTypes":
				return ec.fieldContext___Type_possibleTypes(ctx, field)
			case "enumValues":
				return ec.fieldContext___Type_enumValues(ctx, field)
			case "inputFields":
				return ec.fieldContext___Type_inputFields(ctx, field)
			case "ofType":
				return ec.fieldContext___Type_ofType(ctx, field)
			case "specifiedByURL":
				return ec.fieldContext___Type_specifiedByURL(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type __Type", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) ___InputValue_defaultValue(ctx context.Context, field graphql.CollectedField, obj *introspection.InputValue) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___InputValue_defaultValue(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.DefaultValue, nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.(*string)
	fc.Result = res
	return ec.marshalOString2ᚖstring(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___InputValue_defaultValue(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__InputValue",
		Field:      field,
		IsMethod:   false,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Schema_description(ctx context.Context, field graphql.CollectedField, obj *introspection.Schema) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Schema_description(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Description(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.(*string)
	fc.Result = res
	return ec.marshalOString2ᚖstring(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Schema_description(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Schema",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Schema_types(ctx context.Context, field graphql.CollectedField, obj *introspection.Schema) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Schema_types(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Types(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.([]introspection.Type)
	fc.Result = res
	return ec.marshalN__Type2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐTypeᚄ(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Schema_types(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Schema",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "kind":
				return ec.fieldContext___Type_kind(ctx, field)
			case "name":
				return ec.fieldContext___Type_name(ctx, field)
			case "description":
				return ec.fieldContext___Type_description(ctx, field)
			case "fields":
				return ec.fieldContext___Type_fields(ctx, field)
			case "interfaces":
				return ec.fieldContext___Type_interfaces(ctx, field)
			case "possibleTypes":
				return ec.fieldContext___Type_possibleTypes(ctx, field)
			case "enumValues":
				return ec.fieldContext___Type_enumValues(ctx, field)
			case "inputFields":
				return ec.fieldContext___Type_inputFields(ctx, field)
			case "ofType":
				return ec.fieldContext___Type_ofType(ctx, field)
			case "specifiedByURL":
				return ec.fieldContext___Type_specifiedByURL(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type __Type", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Schema_queryType(ctx context.Context, field graphql.CollectedField, obj *introspection.Schema) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Schema_queryType(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.QueryType(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(*introspection.Type)
	fc.Result = res
	return ec.marshalN__Type2ᚖgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐType(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Schema_queryType(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Schema",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "kind":
				return ec.fieldContext___Type_kind(ctx, field)
			case "name":
				return ec.fieldContext___Type_name(ctx, field)
			case "description":
				return ec.fieldContext___Type_description(ctx, field)
			case "fields":
				return ec.fieldContext___Type_fields(ctx, field)
			case "interfaces":
				return ec.fieldContext___Type_interfaces(ctx, field)
			case "possibleTypes":
				return ec.fieldContext___Type_possibleTypes(ctx, field)
			case "enumValues":
				return ec.fieldContext___Type_enumValues(ctx, field)
			case "inputFields":
				return ec.fieldContext___Type_inputFields(ctx, field)
			case "ofType":
				return ec.fieldContext___Type_ofType(ctx, field)
			case "specifiedByURL":
				return ec.fieldContext___Type_specifiedByURL(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type __Type", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Schema_mutationType(ctx context.Context, field graphql.CollectedField, obj *introspection.Schema) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Schema_mutationType(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.MutationType(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.(*introspection.Type)
	fc.Result = res
	return ec.marshalO__Type2ᚖgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐType(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Schema_mutationType(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Schema",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "kind":
				return ec.fieldContext___Type_kind(ctx, field)
			case "name":
				return ec.fieldContext___Type_name(ctx, field)
			case "description":
				return ec.fieldContext___Type_description(ctx, field)
			case "fields":
				return ec.fieldContext___Type_fields(ctx, field)
			case "interfaces":
				return ec.fieldContext___Type_interfaces(ctx, field)
			case "possibleTypes":
				return ec.fieldContext___Type_possibleTypes(ctx, field)
			case "enumValues":
				return ec.fieldContext___Type_enumValues(ctx, field)
			case "inputFields":
				return ec.fieldContext___Type_inputFields(ctx, field)
			case "ofType":
				return ec.fieldContext___Type_ofType(ctx, field)
			case "specifiedByURL":
				return ec.fieldContext___Type_specifiedByURL(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type __Type", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Schema_subscriptionType(ctx context.Context, field graphql.CollectedField, obj *introspection.Schema) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Schema_subscriptionType(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.SubscriptionType(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.(*introspection.Type)
	fc.Result = res
	return ec.marshalO__Type2ᚖgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐType(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Schema_subscriptionType(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Schema",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "kind":
				return ec.fieldContext___Type_kind(ctx, field)
			case "name":
				return ec.fieldContext___Type_name(ctx, field)
			case "description":
				return ec.fieldContext___Type_description(ctx, field)
			case "fields":
				return ec.fieldContext___Type_fields(ctx, field)
			case "interfaces":
				return ec.fieldContext___Type_interfaces(ctx, field)
			case "possibleTypes":
				return ec.fieldContext___Type_possibleTypes(ctx, field)
			case "enumValues":
				return ec.fieldContext___Type_enumValues(ctx, field)
			case "inputFields":
				return ec.fieldContext___Type_inputFields(ctx, field)
			case "ofType":
				return ec.fieldContext___Type_ofType(ctx, field)
			case "specifiedByURL":
				return ec.fieldContext___Type_specifiedByURL(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type __Type", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Schema_directives(ctx context.Context, field graphql.CollectedField, obj *introspection.Schema) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Schema_directives(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Directives(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.([]introspection.Directive)
	fc.Result = res
	return ec.marshalN__Directive2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐDirectiveᚄ(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Schema_directives(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Schema",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "name":
				return ec.fieldContext___Directive_name(ctx, field)
			case "description":
				return ec.fieldContext___Directive_description(ctx, field)
			case "locations":
				return ec.fieldContext___Directive_locations(ctx, field)
			case "args":
				return ec.fieldContext___Directive_args(ctx, field)
			case "isRepeatable":
				return ec.fieldContext___Directive_isRepeatable(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type __Directive", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Type_kind(ctx context.Context, field graphql.CollectedField, obj *introspection.Type) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Type_kind(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Kind(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		if !graphql.HasFieldError(ctx, fc) {
			ec.Errorf(ctx, "must not be null")
		}
		return graphql.Null
	}
	res := resTmp.(string)
	fc.Result = res
	return ec.marshalN__TypeKind2string(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Type_kind(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Type",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type __TypeKind does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Type_name(ctx context.Context, field graphql.CollectedField, obj *introspection.Type) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Type_name(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Name(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.(*string)
	fc.Result = res
	return ec.marshalOString2ᚖstring(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Type_name(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Type",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Type_description(ctx context.Context, field graphql.CollectedField, obj *introspection.Type) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Type_description(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Description(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.(*string)
	fc.Result = res
	return ec.marshalOString2ᚖstring(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Type_description(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Type",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Type_fields(ctx context.Context, field graphql.CollectedField, obj *introspection.Type) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Type_fields(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Fields(fc.Args["includeDeprecated"].(bool)), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.([]introspection.Field)
	fc.Result = res
	return ec.marshalO__Field2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐFieldᚄ(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Type_fields(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Type",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "name":
				return ec.fieldContext___Field_name(ctx, field)
			case "description":
				return ec.fieldContext___Field_description(ctx, field)
			case "args":
				return ec.fieldContext___Field_args(ctx, field)
			case "type":
				return ec.fieldContext___Field_type(ctx, field)
			case "isDeprecated":
				return ec.fieldContext___Field_isDeprecated(ctx, field)
			case "deprecationReason":
				return ec.fieldContext___Field_deprecationReason(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type __Field", field.Name)
		},
	}
	defer func() {
		if r := recover(); r != nil {
			err = ec.Recover(ctx, r)
			ec.Error(ctx, err)
		}
	}()
	ctx = graphql.WithFieldContext(ctx, fc)
	if fc.Args, err = ec.field___Type_fields_args(ctx, field.ArgumentMap(ec.Variables)); err != nil {
		ec.Error(ctx, err)
		return fc, err
	}
	return fc, nil
}

func (ec *executionContext) ___Type_interfaces(ctx context.Context, field graphql.CollectedField, obj *introspection.Type) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Type_interfaces(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.Interfaces(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.([]introspection.Type)
	fc.Result = res
	return ec.marshalO__Type2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐTypeᚄ(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Type_interfaces(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Type",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "kind":
				return ec.fieldContext___Type_kind(ctx, field)
			case "name":
				return ec.fieldContext___Type_name(ctx, field)
			case "description":
				return ec.fieldContext___Type_description(ctx, field)
			case "fields":
				return ec.fieldContext___Type_fields(ctx, field)
			case "interfaces":
				return ec.fieldContext___Type_interfaces(ctx, field)
			case "possibleTypes":
				return ec.fieldContext___Type_possibleTypes(ctx, field)
			case "enumValues":
				return ec.fieldContext___Type_enumValues(ctx, field)
			case "inputFields":
				return ec.fieldContext___Type_inputFields(ctx, field)
			case "ofType":
				return ec.fieldContext___Type_ofType(ctx, field)
			case "specifiedByURL":
				return ec.fieldContext___Type_specifiedByURL(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type __Type", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Type_possibleTypes(ctx context.Context, field graphql.CollectedField, obj *introspection.Type) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Type_possibleTypes(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.PossibleTypes(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.([]introspection.Type)
	fc.Result = res
	return ec.marshalO__Type2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐTypeᚄ(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Type_possibleTypes(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Type",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "kind":
				return ec.fieldContext___Type_kind(ctx, field)
			case "name":
				return ec.fieldContext___Type_name(ctx, field)
			case "description":
				return ec.fieldContext___Type_description(ctx, field)
			case "fields":
				return ec.fieldContext___Type_fields(ctx, field)
			case "interfaces":
				return ec.fieldContext___Type_interfaces(ctx, field)
			case "possibleTypes":
				return ec.fieldContext___Type_possibleTypes(ctx, field)
			case "enumValues":
				return ec.fieldContext___Type_enumValues(ctx, field)
			case "inputFields":
				return ec.fieldContext___Type_inputFields(ctx, field)
			case "ofType":
				return ec.fieldContext___Type_ofType(ctx, field)
			case "specifiedByURL":
				return ec.fieldContext___Type_specifiedByURL(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type __Type", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Type_enumValues(ctx context.Context, field graphql.CollectedField, obj *introspection.Type) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Type_enumValues(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.EnumValues(fc.Args["includeDeprecated"].(bool)), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.([]introspection.EnumValue)
	fc.Result = res
	return ec.marshalO__EnumValue2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐEnumValueᚄ(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Type_enumValues(ctx context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Type",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "name":
				return ec.fieldContext___EnumValue_name(ctx, field)
			case "description":
				return ec.fieldContext___EnumValue_description(ctx, field)
			case "isDeprecated":
				return ec.fieldContext___EnumValue_isDeprecated(ctx, field)
			case "deprecationReason":
				return ec.fieldContext___EnumValue_deprecationReason(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type __EnumValue", field.Name)
		},
	}
	defer func() {
		if r := recover(); r != nil {
			err = ec.Recover(ctx, r)
			ec.Error(ctx, err)
		}
	}()
	ctx = graphql.WithFieldContext(ctx, fc)
	if fc.Args, err = ec.field___Type_enumValues_args(ctx, field.ArgumentMap(ec.Variables)); err != nil {
		ec.Error(ctx, err)
		return fc, err
	}
	return fc, nil
}

func (ec *executionContext) ___Type_inputFields(ctx context.Context, field graphql.CollectedField, obj *introspection.Type) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Type_inputFields(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.InputFields(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.([]introspection.InputValue)
	fc.Result = res
	return ec.marshalO__InputValue2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐInputValueᚄ(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Type_inputFields(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Type",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "name":
				return ec.fieldContext___InputValue_name(ctx, field)
			case "description":
				return ec.fieldContext___InputValue_description(ctx, field)
			case "type":
				return ec.fieldContext___InputValue_type(ctx, field)
			case "defaultValue":
				return ec.fieldContext___InputValue_defaultValue(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type __InputValue", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Type_ofType(ctx context.Context, field graphql.CollectedField, obj *introspection.Type) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Type_ofType(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.OfType(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.(*introspection.Type)
	fc.Result = res
	return ec.marshalO__Type2ᚖgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐType(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Type_ofType(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Type",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			switch field.Name {
			case "kind":
				return ec.fieldContext___Type_kind(ctx, field)
			case "name":
				return ec.fieldContext___Type_name(ctx, field)
			case "description":
				return ec.fieldContext___Type_description(ctx, field)
			case "fields":
				return ec.fieldContext___Type_fields(ctx, field)
			case "interfaces":
				return ec.fieldContext___Type_interfaces(ctx, field)
			case "possibleTypes":
				return ec.fieldContext___Type_possibleTypes(ctx, field)
			case "enumValues":
				return ec.fieldContext___Type_enumValues(ctx, field)
			case "inputFields":
				return ec.fieldContext___Type_inputFields(ctx, field)
			case "ofType":
				return ec.fieldContext___Type_ofType(ctx, field)
			case "specifiedByURL":
				return ec.fieldContext___Type_specifiedByURL(ctx, field)
			}
			return nil, fmt.Errorf("no field named %q was found under type __Type", field.Name)
		},
	}
	return fc, nil
}

func (ec *executionContext) ___Type_specifiedByURL(ctx context.Context, field graphql.CollectedField, obj *introspection.Type) (ret graphql.Marshaler) {
	fc, err := ec.fieldContext___Type_specifiedByURL(ctx, field)
	if err != nil {
		return graphql.Null
	}
	ctx = graphql.WithFieldContext(ctx, fc)
	defer func() {
		if r := recover(); r != nil {
			ec.Error(ctx, ec.Recover(ctx, r))
			ret = graphql.Null
		}
	}()
	resTmp, err := ec.ResolverMiddleware(ctx, func(rctx context.Context) (interface{}, error) {
		ctx = rctx // use context from middleware stack in children
		return obj.SpecifiedByURL(), nil
	})
	if err != nil {
		ec.Error(ctx, err)
		return graphql.Null
	}
	if resTmp == nil {
		return graphql.Null
	}
	res := resTmp.(*string)
	fc.Result = res
	return ec.marshalOString2ᚖstring(ctx, field.Selections, res)
}

func (ec *executionContext) fieldContext___Type_specifiedByURL(_ context.Context, field graphql.CollectedField) (fc *graphql.FieldContext, err error) {
	fc = &graphql.FieldContext{
		Object:     "__Type",
		Field:      field,
		IsMethod:   true,
		IsResolver: false,
		Child: func(ctx context.Context, field graphql.CollectedField) (*graphql.FieldContext, error) {
			return nil, errors.New("field of type String does not have child fields")
		},
	}
	return fc, nil
}

func (ec *executionContext) unmarshalInputCommentFilter(ctx context.Context, obj interface{}) (model.CommentFilter, error) {
	var it model.CommentFilter
	asMap := map[string]interface{}{}
	for k, v := range obj.(map[string]interface{}) {
		asMap[k] = v
	}

	fieldsInOrder := [...]string{"author", "body"}
	for _, k := range fieldsInOrder {
		v, ok := asMap[k]
		if !ok {
			continue
		}
		switch k {
		case "author":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("author"))
			data, err := ec.unmarshalOString2ᚖstring(ctx, v)
			if err != nil {
				return it, err
			}
			it.Author = data
		case "body":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("body"))
			data, err := ec.unmarshalOString2ᚖstring(ctx, v)
			if err != nil {
				return it, err
			}
			it.Body = data
		}
	}

	return it, nil
}

func (ec *executionContext) unmarshalInputNewComment(ctx context.Context, obj interface{}) (model.NewComment, error) {
	var it model.NewComment
	asMap := map[string]interface{}{}
	for k, v := range obj.(map[string]interface{}) {
		asMap[k] = v
	}

	fieldsInOrder := [...]string{"author", "body", "parentPost", "parentComment"}
	for _, k := range fieldsInOrder {
		v, ok := asMap[k]
		if !ok {
			continue
		}
		switch k {
		case "author":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("author"))
			data, err := ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.Author = data
		case "body":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("body"))
			data, err := ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.Body = data
		case "parentPost":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("parentPost"))
			data, err := ec.unmarshalNID2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.ParentPost = data
		case "parentComment":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("parentComment"))
			data, err := ec.unmarshalOID2ᚖstring(ctx, v)
			if err != nil {
				return it, err
			}
			it.ParentComment = data
		}
	}

	return it, nil
}

func (ec *executionContext) unmarshalInputNewPost(ctx context.Context, obj interface{}) (model.NewPost, error) {
	var it model.NewPost
	asMap := map[string]interface{}{}
	for k, v := range obj.(map[string]interface{}) {
		asMap[k] = v
	}

	fieldsInOrder := [...]string{"author", "body", "open"}
	for _, k := range fieldsInOrder {
		v, ok := asMap[k]
		if !ok {
			continue
		}
		switch k {
		case "author":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("author"))
			data, err := ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.Author = data
		case "body":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("body"))
			data, err := ec.unmarshalNString2string(ctx, v)
			if err != nil {
				return it, err
			}
			it.Body = data
		case "open":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("open"))
			data, err := ec.unmarshalNBoolean2bool(ctx, v)
			if err != nil {
				return it, err
			}
			it.Open = data
		}
	}

	return it, nil
}

func (ec *executionContext) unmarshalInputUpdateComment(ctx context.Context, obj interface{}) (model.UpdateComment, error) {
	var it model.UpdateComment
	asMap := map[string]interface{}{}
	for k, v := range obj.(map[string]interface{}) {
		asMap[k] = v
	}

	fieldsInOrder := [...]string{"body"}
	for _, k := range fieldsInOrder {
		v, ok := asMap[k]
		if !ok {
			continue
		}
		switch k {
		case "body":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("body"))
			data, err := ec.unmarshalOString2ᚖstring(ctx, v)
			if err != nil {
				return it, err
			}
			it.Body = data
		}
	}

	return it, nil
}

func (ec *executionContext) unmarshalInputUpdatePost(ctx context.Context, obj interface{}) (model.UpdatePost, error) {
	var it model.UpdatePost
	asMap := map[string]interface{}{}
	for k, v := range obj.(map[string]interface{}) {
		asMap[k] = v
	}

	fieldsInOrder := [...]string{"body", "open"}
	for _, k := range fieldsInOrder {
		v, ok := asMap[k]
		if !ok {
			continue
		}
		switch k {
		case "body":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("body"))
			data, err := ec.unmarshalOString2ᚖstring(ctx, v)
			if err != nil {
				return it, err
			}
			it.Body = data
		case "open":
			ctx := graphql.WithPathContext(ctx, graphql.NewPathWithField("open"))
			data, err := ec.unmarshalOBoolean2ᚖbool(ctx, v)
			if err != nil {
				return it, err
			}
			it.Open = data
		}
	}

	return it, nil
}

var commentImplementors = []string{"Comment"}

func (ec *executionContext) _Comment(ctx context.Context, sel ast.SelectionSet, obj *models.Comment) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, commentImplementors)

	out := graphql.NewFieldSet(fields)
	deferred := make(map[string]*graphql.FieldSet)
	for i, field := range fields {
		switch field.Name {
		case "__typename":
			out.Values[i] = graphql.MarshalString("Comment")
		case "id":
			out.Values[i] = ec._Comment_id(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				atomic.AddUint32(&out.Invalids, 1)
			}
		case "author":
			out.Values[i] = ec._Comment_author(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				atomic.AddUint32(&out.Invalids, 1)
			}
		case "body":
			out.Values[i] = ec._Comment_body(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				atomic.AddUint32(&out.Invalids, 1)
			}
		case "comments":
			field := field

			innerFunc := func(ctx context.Context, _ *graphql.FieldSet) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Comment_comments(ctx, field, obj)
				return res
			}

			if field.Deferrable != nil {
				dfs, ok := deferred[field.Deferrable.Label]
				di := 0
				if ok {
					dfs.AddField(field)
					di = len(dfs.Values) - 1
				} else {
					dfs = graphql.NewFieldSet([]graphql.CollectedField{field})
					deferred[field.Deferrable.Label] = dfs
				}
				dfs.Concurrently(di, func(ctx context.Context) graphql.Marshaler {
					return innerFunc(ctx, dfs)
				})

				// don't run the out.Concurrently() call below
				out.Values[i] = graphql.Null
				continue
			}

			out.Concurrently(i, func(ctx context.Context) graphql.Marshaler { return innerFunc(ctx, out) })
		default:
			panic("unknown field " + strconv.Quote(field.Name))
		}
	}
	out.Dispatch(ctx)
	if out.Invalids > 0 {
		return graphql.Null
	}

	atomic.AddInt32(&ec.deferred, int32(len(deferred)))

	for label, dfs := range deferred {
		ec.processDeferredGroup(graphql.DeferredGroup{
			Label:    label,
			Path:     graphql.GetPath(ctx),
			FieldSet: dfs,
			Context:  ctx,
		})
	}

	return out
}

var mutationImplementors = []string{"Mutation"}

func (ec *executionContext) _Mutation(ctx context.Context, sel ast.SelectionSet) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, mutationImplementors)
	ctx = graphql.WithFieldContext(ctx, &graphql.FieldContext{
		Object: "Mutation",
	})

	out := graphql.NewFieldSet(fields)
	deferred := make(map[string]*graphql.FieldSet)
	for i, field := range fields {
		innerCtx := graphql.WithRootFieldContext(ctx, &graphql.RootFieldContext{
			Object: field.Name,
			Field:  field,
		})

		switch field.Name {
		case "__typename":
			out.Values[i] = graphql.MarshalString("Mutation")
		case "createComment":
			out.Values[i] = ec.OperationContext.RootResolverMiddleware(innerCtx, func(ctx context.Context) (res graphql.Marshaler) {
				return ec._Mutation_createComment(ctx, field)
			})
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "updateComment":
			out.Values[i] = ec.OperationContext.RootResolverMiddleware(innerCtx, func(ctx context.Context) (res graphql.Marshaler) {
				return ec._Mutation_updateComment(ctx, field)
			})
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "deleteComment":
			out.Values[i] = ec.OperationContext.RootResolverMiddleware(innerCtx, func(ctx context.Context) (res graphql.Marshaler) {
				return ec._Mutation_deleteComment(ctx, field)
			})
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "createPost":
			out.Values[i] = ec.OperationContext.RootResolverMiddleware(innerCtx, func(ctx context.Context) (res graphql.Marshaler) {
				return ec._Mutation_createPost(ctx, field)
			})
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "updatePost":
			out.Values[i] = ec.OperationContext.RootResolverMiddleware(innerCtx, func(ctx context.Context) (res graphql.Marshaler) {
				return ec._Mutation_updatePost(ctx, field)
			})
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "deletePost":
			out.Values[i] = ec.OperationContext.RootResolverMiddleware(innerCtx, func(ctx context.Context) (res graphql.Marshaler) {
				return ec._Mutation_deletePost(ctx, field)
			})
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		default:
			panic("unknown field " + strconv.Quote(field.Name))
		}
	}
	out.Dispatch(ctx)
	if out.Invalids > 0 {
		return graphql.Null
	}

	atomic.AddInt32(&ec.deferred, int32(len(deferred)))

	for label, dfs := range deferred {
		ec.processDeferredGroup(graphql.DeferredGroup{
			Label:    label,
			Path:     graphql.GetPath(ctx),
			FieldSet: dfs,
			Context:  ctx,
		})
	}

	return out
}

var postImplementors = []string{"Post"}

func (ec *executionContext) _Post(ctx context.Context, sel ast.SelectionSet, obj *models.Post) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, postImplementors)

	out := graphql.NewFieldSet(fields)
	deferred := make(map[string]*graphql.FieldSet)
	for i, field := range fields {
		switch field.Name {
		case "__typename":
			out.Values[i] = graphql.MarshalString("Post")
		case "id":
			out.Values[i] = ec._Post_id(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				atomic.AddUint32(&out.Invalids, 1)
			}
		case "author":
			out.Values[i] = ec._Post_author(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				atomic.AddUint32(&out.Invalids, 1)
			}
		case "body":
			out.Values[i] = ec._Post_body(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				atomic.AddUint32(&out.Invalids, 1)
			}
		case "open":
			out.Values[i] = ec._Post_open(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				atomic.AddUint32(&out.Invalids, 1)
			}
		case "comments":
			field := field

			innerFunc := func(ctx context.Context, _ *graphql.FieldSet) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Post_comments(ctx, field, obj)
				return res
			}

			if field.Deferrable != nil {
				dfs, ok := deferred[field.Deferrable.Label]
				di := 0
				if ok {
					dfs.AddField(field)
					di = len(dfs.Values) - 1
				} else {
					dfs = graphql.NewFieldSet([]graphql.CollectedField{field})
					deferred[field.Deferrable.Label] = dfs
				}
				dfs.Concurrently(di, func(ctx context.Context) graphql.Marshaler {
					return innerFunc(ctx, dfs)
				})

				// don't run the out.Concurrently() call below
				out.Values[i] = graphql.Null
				continue
			}

			out.Concurrently(i, func(ctx context.Context) graphql.Marshaler { return innerFunc(ctx, out) })
		default:
			panic("unknown field " + strconv.Quote(field.Name))
		}
	}
	out.Dispatch(ctx)
	if out.Invalids > 0 {
		return graphql.Null
	}

	atomic.AddInt32(&ec.deferred, int32(len(deferred)))

	for label, dfs := range deferred {
		ec.processDeferredGroup(graphql.DeferredGroup{
			Label:    label,
			Path:     graphql.GetPath(ctx),
			FieldSet: dfs,
			Context:  ctx,
		})
	}

	return out
}

var queryImplementors = []string{"Query"}

func (ec *executionContext) _Query(ctx context.Context, sel ast.SelectionSet) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, queryImplementors)
	ctx = graphql.WithFieldContext(ctx, &graphql.FieldContext{
		Object: "Query",
	})

	out := graphql.NewFieldSet(fields)
	deferred := make(map[string]*graphql.FieldSet)
	for i, field := range fields {
		innerCtx := graphql.WithRootFieldContext(ctx, &graphql.RootFieldContext{
			Object: field.Name,
			Field:  field,
		})

		switch field.Name {
		case "__typename":
			out.Values[i] = graphql.MarshalString("Query")
		case "comments":
			field := field

			innerFunc := func(ctx context.Context, fs *graphql.FieldSet) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Query_comments(ctx, field)
				if res == graphql.Null {
					atomic.AddUint32(&fs.Invalids, 1)
				}
				return res
			}

			rrm := func(ctx context.Context) graphql.Marshaler {
				return ec.OperationContext.RootResolverMiddleware(ctx,
					func(ctx context.Context) graphql.Marshaler { return innerFunc(ctx, out) })
			}

			out.Concurrently(i, func(ctx context.Context) graphql.Marshaler { return rrm(innerCtx) })
		case "posts":
			field := field

			innerFunc := func(ctx context.Context, fs *graphql.FieldSet) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Query_posts(ctx, field)
				if res == graphql.Null {
					atomic.AddUint32(&fs.Invalids, 1)
				}
				return res
			}

			rrm := func(ctx context.Context) graphql.Marshaler {
				return ec.OperationContext.RootResolverMiddleware(ctx,
					func(ctx context.Context) graphql.Marshaler { return innerFunc(ctx, out) })
			}

			out.Concurrently(i, func(ctx context.Context) graphql.Marshaler { return rrm(innerCtx) })
		case "post":
			field := field

			innerFunc := func(ctx context.Context, fs *graphql.FieldSet) (res graphql.Marshaler) {
				defer func() {
					if r := recover(); r != nil {
						ec.Error(ctx, ec.Recover(ctx, r))
					}
				}()
				res = ec._Query_post(ctx, field)
				if res == graphql.Null {
					atomic.AddUint32(&fs.Invalids, 1)
				}
				return res
			}

			rrm := func(ctx context.Context) graphql.Marshaler {
				return ec.OperationContext.RootResolverMiddleware(ctx,
					func(ctx context.Context) graphql.Marshaler { return innerFunc(ctx, out) })
			}

			out.Concurrently(i, func(ctx context.Context) graphql.Marshaler { return rrm(innerCtx) })
		case "__type":
			out.Values[i] = ec.OperationContext.RootResolverMiddleware(innerCtx, func(ctx context.Context) (res graphql.Marshaler) {
				return ec._Query___type(ctx, field)
			})
		case "__schema":
			out.Values[i] = ec.OperationContext.RootResolverMiddleware(innerCtx, func(ctx context.Context) (res graphql.Marshaler) {
				return ec._Query___schema(ctx, field)
			})
		default:
			panic("unknown field " + strconv.Quote(field.Name))
		}
	}
	out.Dispatch(ctx)
	if out.Invalids > 0 {
		return graphql.Null
	}

	atomic.AddInt32(&ec.deferred, int32(len(deferred)))

	for label, dfs := range deferred {
		ec.processDeferredGroup(graphql.DeferredGroup{
			Label:    label,
			Path:     graphql.GetPath(ctx),
			FieldSet: dfs,
			Context:  ctx,
		})
	}

	return out
}

var subscriptionImplementors = []string{"Subscription"}

func (ec *executionContext) _Subscription(ctx context.Context, sel ast.SelectionSet) func(ctx context.Context) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, subscriptionImplementors)
	ctx = graphql.WithFieldContext(ctx, &graphql.FieldContext{
		Object: "Subscription",
	})
	if len(fields) != 1 {
		ec.Errorf(ctx, "must subscribe to exactly one stream")
		return nil
	}

	switch fields[0].Name {
	case "commentAdded":
		return ec._Subscription_commentAdded(ctx, fields[0])
	default:
		panic("unknown field " + strconv.Quote(fields[0].Name))
	}
}

var __DirectiveImplementors = []string{"__Directive"}

func (ec *executionContext) ___Directive(ctx context.Context, sel ast.SelectionSet, obj *introspection.Directive) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, __DirectiveImplementors)

	out := graphql.NewFieldSet(fields)
	deferred := make(map[string]*graphql.FieldSet)
	for i, field := range fields {
		switch field.Name {
		case "__typename":
			out.Values[i] = graphql.MarshalString("__Directive")
		case "name":
			out.Values[i] = ec.___Directive_name(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "description":
			out.Values[i] = ec.___Directive_description(ctx, field, obj)
		case "locations":
			out.Values[i] = ec.___Directive_locations(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "args":
			out.Values[i] = ec.___Directive_args(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "isRepeatable":
			out.Values[i] = ec.___Directive_isRepeatable(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		default:
			panic("unknown field " + strconv.Quote(field.Name))
		}
	}
	out.Dispatch(ctx)
	if out.Invalids > 0 {
		return graphql.Null
	}

	atomic.AddInt32(&ec.deferred, int32(len(deferred)))

	for label, dfs := range deferred {
		ec.processDeferredGroup(graphql.DeferredGroup{
			Label:    label,
			Path:     graphql.GetPath(ctx),
			FieldSet: dfs,
			Context:  ctx,
		})
	}

	return out
}

var __EnumValueImplementors = []string{"__EnumValue"}

func (ec *executionContext) ___EnumValue(ctx context.Context, sel ast.SelectionSet, obj *introspection.EnumValue) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, __EnumValueImplementors)

	out := graphql.NewFieldSet(fields)
	deferred := make(map[string]*graphql.FieldSet)
	for i, field := range fields {
		switch field.Name {
		case "__typename":
			out.Values[i] = graphql.MarshalString("__EnumValue")
		case "name":
			out.Values[i] = ec.___EnumValue_name(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "description":
			out.Values[i] = ec.___EnumValue_description(ctx, field, obj)
		case "isDeprecated":
			out.Values[i] = ec.___EnumValue_isDeprecated(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "deprecationReason":
			out.Values[i] = ec.___EnumValue_deprecationReason(ctx, field, obj)
		default:
			panic("unknown field " + strconv.Quote(field.Name))
		}
	}
	out.Dispatch(ctx)
	if out.Invalids > 0 {
		return graphql.Null
	}

	atomic.AddInt32(&ec.deferred, int32(len(deferred)))

	for label, dfs := range deferred {
		ec.processDeferredGroup(graphql.DeferredGroup{
			Label:    label,
			Path:     graphql.GetPath(ctx),
			FieldSet: dfs,
			Context:  ctx,
		})
	}

	return out
}

var __FieldImplementors = []string{"__Field"}

func (ec *executionContext) ___Field(ctx context.Context, sel ast.SelectionSet, obj *introspection.Field) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, __FieldImplementors)

	out := graphql.NewFieldSet(fields)
	deferred := make(map[string]*graphql.FieldSet)
	for i, field := range fields {
		switch field.Name {
		case "__typename":
			out.Values[i] = graphql.MarshalString("__Field")
		case "name":
			out.Values[i] = ec.___Field_name(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "description":
			out.Values[i] = ec.___Field_description(ctx, field, obj)
		case "args":
			out.Values[i] = ec.___Field_args(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "type":
			out.Values[i] = ec.___Field_type(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "isDeprecated":
			out.Values[i] = ec.___Field_isDeprecated(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "deprecationReason":
			out.Values[i] = ec.___Field_deprecationReason(ctx, field, obj)
		default:
			panic("unknown field " + strconv.Quote(field.Name))
		}
	}
	out.Dispatch(ctx)
	if out.Invalids > 0 {
		return graphql.Null
	}

	atomic.AddInt32(&ec.deferred, int32(len(deferred)))

	for label, dfs := range deferred {
		ec.processDeferredGroup(graphql.DeferredGroup{
			Label:    label,
			Path:     graphql.GetPath(ctx),
			FieldSet: dfs,
			Context:  ctx,
		})
	}

	return out
}

var __InputValueImplementors = []string{"__InputValue"}

func (ec *executionContext) ___InputValue(ctx context.Context, sel ast.SelectionSet, obj *introspection.InputValue) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, __InputValueImplementors)

	out := graphql.NewFieldSet(fields)
	deferred := make(map[string]*graphql.FieldSet)
	for i, field := range fields {
		switch field.Name {
		case "__typename":
			out.Values[i] = graphql.MarshalString("__InputValue")
		case "name":
			out.Values[i] = ec.___InputValue_name(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "description":
			out.Values[i] = ec.___InputValue_description(ctx, field, obj)
		case "type":
			out.Values[i] = ec.___InputValue_type(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "defaultValue":
			out.Values[i] = ec.___InputValue_defaultValue(ctx, field, obj)
		default:
			panic("unknown field " + strconv.Quote(field.Name))
		}
	}
	out.Dispatch(ctx)
	if out.Invalids > 0 {
		return graphql.Null
	}

	atomic.AddInt32(&ec.deferred, int32(len(deferred)))

	for label, dfs := range deferred {
		ec.processDeferredGroup(graphql.DeferredGroup{
			Label:    label,
			Path:     graphql.GetPath(ctx),
			FieldSet: dfs,
			Context:  ctx,
		})
	}

	return out
}

var __SchemaImplementors = []string{"__Schema"}

func (ec *executionContext) ___Schema(ctx context.Context, sel ast.SelectionSet, obj *introspection.Schema) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, __SchemaImplementors)

	out := graphql.NewFieldSet(fields)
	deferred := make(map[string]*graphql.FieldSet)
	for i, field := range fields {
		switch field.Name {
		case "__typename":
			out.Values[i] = graphql.MarshalString("__Schema")
		case "description":
			out.Values[i] = ec.___Schema_description(ctx, field, obj)
		case "types":
			out.Values[i] = ec.___Schema_types(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "queryType":
			out.Values[i] = ec.___Schema_queryType(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "mutationType":
			out.Values[i] = ec.___Schema_mutationType(ctx, field, obj)
		case "subscriptionType":
			out.Values[i] = ec.___Schema_subscriptionType(ctx, field, obj)
		case "directives":
			out.Values[i] = ec.___Schema_directives(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		default:
			panic("unknown field " + strconv.Quote(field.Name))
		}
	}
	out.Dispatch(ctx)
	if out.Invalids > 0 {
		return graphql.Null
	}

	atomic.AddInt32(&ec.deferred, int32(len(deferred)))

	for label, dfs := range deferred {
		ec.processDeferredGroup(graphql.DeferredGroup{
			Label:    label,
			Path:     graphql.GetPath(ctx),
			FieldSet: dfs,
			Context:  ctx,
		})
	}

	return out
}

var __TypeImplementors = []string{"__Type"}

func (ec *executionContext) ___Type(ctx context.Context, sel ast.SelectionSet, obj *introspection.Type) graphql.Marshaler {
	fields := graphql.CollectFields(ec.OperationContext, sel, __TypeImplementors)

	out := graphql.NewFieldSet(fields)
	deferred := make(map[string]*graphql.FieldSet)
	for i, field := range fields {
		switch field.Name {
		case "__typename":
			out.Values[i] = graphql.MarshalString("__Type")
		case "kind":
			out.Values[i] = ec.___Type_kind(ctx, field, obj)
			if out.Values[i] == graphql.Null {
				out.Invalids++
			}
		case "name":
			out.Values[i] = ec.___Type_name(ctx, field, obj)
		case "description":
			out.Values[i] = ec.___Type_description(ctx, field, obj)
		case "fields":
			out.Values[i] = ec.___Type_fields(ctx, field, obj)
		case "interfaces":
			out.Values[i] = ec.___Type_interfaces(ctx, field, obj)
		case "possibleTypes":
			out.Values[i] = ec.___Type_possibleTypes(ctx, field, obj)
		case "enumValues":
			out.Values[i] = ec.___Type_enumValues(ctx, field, obj)
		case "inputFields":
			out.Values[i] = ec.___Type_inputFields(ctx, field, obj)
		case "ofType":
			out.Values[i] = ec.___Type_ofType(ctx, field, obj)
		case "specifiedByURL":
			out.Values[i] = ec.___Type_specifiedByURL(ctx, field, obj)
		default:
			panic("unknown field " + strconv.Quote(field.Name))
		}
	}
	out.Dispatch(ctx)
	if out.Invalids > 0 {
		return graphql.Null
	}

	atomic.AddInt32(&ec.deferred, int32(len(deferred)))

	for label, dfs := range deferred {
		ec.processDeferredGroup(graphql.DeferredGroup{
			Label:    label,
			Path:     graphql.GetPath(ctx),
			FieldSet: dfs,
			Context:  ctx,
		})
	}

	return out
}

func (ec *executionContext) unmarshalNBoolean2bool(ctx context.Context, v interface{}) (bool, error) {
	res, err := graphql.UnmarshalBoolean(v)
	return res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalNBoolean2bool(ctx context.Context, sel ast.SelectionSet, v bool) graphql.Marshaler {
	res := graphql.MarshalBoolean(v)
	if res == graphql.Null {
		if !graphql.HasFieldError(ctx, graphql.GetFieldContext(ctx)) {
			ec.Errorf(ctx, "the requested element is null which the schema does not allow")
		}
	}
	return res
}

func (ec *executionContext) marshalNComment2githubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐComment(ctx context.Context, sel ast.SelectionSet, v models.Comment) graphql.Marshaler {
	return ec._Comment(ctx, sel, &v)
}

func (ec *executionContext) marshalNComment2ᚕᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐCommentᚄ(ctx context.Context, sel ast.SelectionSet, v []*models.Comment) graphql.Marshaler {
	ret := make(graphql.Array, len(v))
	var wg sync.WaitGroup
	isLen1 := len(v) == 1
	if !isLen1 {
		wg.Add(len(v))
	}
	for i := range v {
		i := i
		fc := &graphql.FieldContext{
			Index:  &i,
			Result: &v[i],
		}
		ctx := graphql.WithFieldContext(ctx, fc)
		f := func(i int) {
			defer func() {
				if r := recover(); r != nil {
					ec.Error(ctx, ec.Recover(ctx, r))
					ret = nil
				}
			}()
			if !isLen1 {
				defer wg.Done()
			}
			ret[i] = ec.marshalNComment2ᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐComment(ctx, sel, v[i])
		}
		if isLen1 {
			f(i)
		} else {
			go f(i)
		}

	}
	wg.Wait()

	for _, e := range ret {
		if e == graphql.Null {
			return graphql.Null
		}
	}

	return ret
}

func (ec *executionContext) marshalNComment2ᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐComment(ctx context.Context, sel ast.SelectionSet, v *models.Comment) graphql.Marshaler {
	if v == nil {
		if !graphql.HasFieldError(ctx, graphql.GetFieldContext(ctx)) {
			ec.Errorf(ctx, "the requested element is null which the schema does not allow")
		}
		return graphql.Null
	}
	return ec._Comment(ctx, sel, v)
}

func (ec *executionContext) unmarshalNID2string(ctx context.Context, v interface{}) (string, error) {
	res, err := graphql.UnmarshalID(v)
	return res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalNID2string(ctx context.Context, sel ast.SelectionSet, v string) graphql.Marshaler {
	res := graphql.MarshalID(v)
	if res == graphql.Null {
		if !graphql.HasFieldError(ctx, graphql.GetFieldContext(ctx)) {
			ec.Errorf(ctx, "the requested element is null which the schema does not allow")
		}
	}
	return res
}

func (ec *executionContext) unmarshalNNewComment2githubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋgraphᚋmodelᚐNewComment(ctx context.Context, v interface{}) (model.NewComment, error) {
	res, err := ec.unmarshalInputNewComment(ctx, v)
	return res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) unmarshalNNewPost2githubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋgraphᚋmodelᚐNewPost(ctx context.Context, v interface{}) (model.NewPost, error) {
	res, err := ec.unmarshalInputNewPost(ctx, v)
	return res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalNPost2githubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐPost(ctx context.Context, sel ast.SelectionSet, v models.Post) graphql.Marshaler {
	return ec._Post(ctx, sel, &v)
}

func (ec *executionContext) marshalNPost2ᚕᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐPostᚄ(ctx context.Context, sel ast.SelectionSet, v []*models.Post) graphql.Marshaler {
	ret := make(graphql.Array, len(v))
	var wg sync.WaitGroup
	isLen1 := len(v) == 1
	if !isLen1 {
		wg.Add(len(v))
	}
	for i := range v {
		i := i
		fc := &graphql.FieldContext{
			Index:  &i,
			Result: &v[i],
		}
		ctx := graphql.WithFieldContext(ctx, fc)
		f := func(i int) {
			defer func() {
				if r := recover(); r != nil {
					ec.Error(ctx, ec.Recover(ctx, r))
					ret = nil
				}
			}()
			if !isLen1 {
				defer wg.Done()
			}
			ret[i] = ec.marshalNPost2ᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐPost(ctx, sel, v[i])
		}
		if isLen1 {
			f(i)
		} else {
			go f(i)
		}

	}
	wg.Wait()

	for _, e := range ret {
		if e == graphql.Null {
			return graphql.Null
		}
	}

	return ret
}

func (ec *executionContext) marshalNPost2ᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐPost(ctx context.Context, sel ast.SelectionSet, v *models.Post) graphql.Marshaler {
	if v == nil {
		if !graphql.HasFieldError(ctx, graphql.GetFieldContext(ctx)) {
			ec.Errorf(ctx, "the requested element is null which the schema does not allow")
		}
		return graphql.Null
	}
	return ec._Post(ctx, sel, v)
}

func (ec *executionContext) unmarshalNString2string(ctx context.Context, v interface{}) (string, error) {
	res, err := graphql.UnmarshalString(v)
	return res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalNString2string(ctx context.Context, sel ast.SelectionSet, v string) graphql.Marshaler {
	res := graphql.MarshalString(v)
	if res == graphql.Null {
		if !graphql.HasFieldError(ctx, graphql.GetFieldContext(ctx)) {
			ec.Errorf(ctx, "the requested element is null which the schema does not allow")
		}
	}
	return res
}

func (ec *executionContext) unmarshalNUpdateComment2githubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋgraphᚋmodelᚐUpdateComment(ctx context.Context, v interface{}) (model.UpdateComment, error) {
	res, err := ec.unmarshalInputUpdateComment(ctx, v)
	return res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) unmarshalNUpdatePost2githubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋgraphᚋmodelᚐUpdatePost(ctx context.Context, v interface{}) (model.UpdatePost, error) {
	res, err := ec.unmarshalInputUpdatePost(ctx, v)
	return res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalN__Directive2githubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐDirective(ctx context.Context, sel ast.SelectionSet, v introspection.Directive) graphql.Marshaler {
	return ec.___Directive(ctx, sel, &v)
}

func (ec *executionContext) marshalN__Directive2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐDirectiveᚄ(ctx context.Context, sel ast.SelectionSet, v []introspection.Directive) graphql.Marshaler {
	ret := make(graphql.Array, len(v))
	var wg sync.WaitGroup
	isLen1 := len(v) == 1
	if !isLen1 {
		wg.Add(len(v))
	}
	for i := range v {
		i := i
		fc := &graphql.FieldContext{
			Index:  &i,
			Result: &v[i],
		}
		ctx := graphql.WithFieldContext(ctx, fc)
		f := func(i int) {
			defer func() {
				if r := recover(); r != nil {
					ec.Error(ctx, ec.Recover(ctx, r))
					ret = nil
				}
			}()
			if !isLen1 {
				defer wg.Done()
			}
			ret[i] = ec.marshalN__Directive2githubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐDirective(ctx, sel, v[i])
		}
		if isLen1 {
			f(i)
		} else {
			go f(i)
		}

	}
	wg.Wait()

	for _, e := range ret {
		if e == graphql.Null {
			return graphql.Null
		}
	}

	return ret
}

func (ec *executionContext) unmarshalN__DirectiveLocation2string(ctx context.Context, v interface{}) (string, error) {
	res, err := graphql.UnmarshalString(v)
	return res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalN__DirectiveLocation2string(ctx context.Context, sel ast.SelectionSet, v string) graphql.Marshaler {
	res := graphql.MarshalString(v)
	if res == graphql.Null {
		if !graphql.HasFieldError(ctx, graphql.GetFieldContext(ctx)) {
			ec.Errorf(ctx, "the requested element is null which the schema does not allow")
		}
	}
	return res
}

func (ec *executionContext) unmarshalN__DirectiveLocation2ᚕstringᚄ(ctx context.Context, v interface{}) ([]string, error) {
	var vSlice []interface{}
	if v != nil {
		vSlice = graphql.CoerceList(v)
	}
	var err error
	res := make([]string, len(vSlice))
	for i := range vSlice {
		ctx := graphql.WithPathContext(ctx, graphql.NewPathWithIndex(i))
		res[i], err = ec.unmarshalN__DirectiveLocation2string(ctx, vSlice[i])
		if err != nil {
			return nil, err
		}
	}
	return res, nil
}

func (ec *executionContext) marshalN__DirectiveLocation2ᚕstringᚄ(ctx context.Context, sel ast.SelectionSet, v []string) graphql.Marshaler {
	ret := make(graphql.Array, len(v))
	var wg sync.WaitGroup
	isLen1 := len(v) == 1
	if !isLen1 {
		wg.Add(len(v))
	}
	for i := range v {
		i := i
		fc := &graphql.FieldContext{
			Index:  &i,
			Result: &v[i],
		}
		ctx := graphql.WithFieldContext(ctx, fc)
		f := func(i int) {
			defer func() {
				if r := recover(); r != nil {
					ec.Error(ctx, ec.Recover(ctx, r))
					ret = nil
				}
			}()
			if !isLen1 {
				defer wg.Done()
			}
			ret[i] = ec.marshalN__DirectiveLocation2string(ctx, sel, v[i])
		}
		if isLen1 {
			f(i)
		} else {
			go f(i)
		}

	}
	wg.Wait()

	for _, e := range ret {
		if e == graphql.Null {
			return graphql.Null
		}
	}

	return ret
}

func (ec *executionContext) marshalN__EnumValue2githubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐEnumValue(ctx context.Context, sel ast.SelectionSet, v introspection.EnumValue) graphql.Marshaler {
	return ec.___EnumValue(ctx, sel, &v)
}

func (ec *executionContext) marshalN__Field2githubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐField(ctx context.Context, sel ast.SelectionSet, v introspection.Field) graphql.Marshaler {
	return ec.___Field(ctx, sel, &v)
}

func (ec *executionContext) marshalN__InputValue2githubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐInputValue(ctx context.Context, sel ast.SelectionSet, v introspection.InputValue) graphql.Marshaler {
	return ec.___InputValue(ctx, sel, &v)
}

func (ec *executionContext) marshalN__InputValue2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐInputValueᚄ(ctx context.Context, sel ast.SelectionSet, v []introspection.InputValue) graphql.Marshaler {
	ret := make(graphql.Array, len(v))
	var wg sync.WaitGroup
	isLen1 := len(v) == 1
	if !isLen1 {
		wg.Add(len(v))
	}
	for i := range v {
		i := i
		fc := &graphql.FieldContext{
			Index:  &i,
			Result: &v[i],
		}
		ctx := graphql.WithFieldContext(ctx, fc)
		f := func(i int) {
			defer func() {
				if r := recover(); r != nil {
					ec.Error(ctx, ec.Recover(ctx, r))
					ret = nil
				}
			}()
			if !isLen1 {
				defer wg.Done()
			}
			ret[i] = ec.marshalN__InputValue2githubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐInputValue(ctx, sel, v[i])
		}
		if isLen1 {
			f(i)
		} else {
			go f(i)
		}

	}
	wg.Wait()

	for _, e := range ret {
		if e == graphql.Null {
			return graphql.Null
		}
	}

	return ret
}

func (ec *executionContext) marshalN__Type2githubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐType(ctx context.Context, sel ast.SelectionSet, v introspection.Type) graphql.Marshaler {
	return ec.___Type(ctx, sel, &v)
}

func (ec *executionContext) marshalN__Type2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐTypeᚄ(ctx context.Context, sel ast.SelectionSet, v []introspection.Type) graphql.Marshaler {
	ret := make(graphql.Array, len(v))
	var wg sync.WaitGroup
	isLen1 := len(v) == 1
	if !isLen1 {
		wg.Add(len(v))
	}
	for i := range v {
		i := i
		fc := &graphql.FieldContext{
			Index:  &i,
			Result: &v[i],
		}
		ctx := graphql.WithFieldContext(ctx, fc)
		f := func(i int) {
			defer func() {
				if r := recover(); r != nil {
					ec.Error(ctx, ec.Recover(ctx, r))
					ret = nil
				}
			}()
			if !isLen1 {
				defer wg.Done()
			}
			ret[i] = ec.marshalN__Type2githubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐType(ctx, sel, v[i])
		}
		if isLen1 {
			f(i)
		} else {
			go f(i)
		}

	}
	wg.Wait()

	for _, e := range ret {
		if e == graphql.Null {
			return graphql.Null
		}
	}

	return ret
}

func (ec *executionContext) marshalN__Type2ᚖgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐType(ctx context.Context, sel ast.SelectionSet, v *introspection.Type) graphql.Marshaler {
	if v == nil {
		if !graphql.HasFieldError(ctx, graphql.GetFieldContext(ctx)) {
			ec.Errorf(ctx, "the requested element is null which the schema does not allow")
		}
		return graphql.Null
	}
	return ec.___Type(ctx, sel, v)
}

func (ec *executionContext) unmarshalN__TypeKind2string(ctx context.Context, v interface{}) (string, error) {
	res, err := graphql.UnmarshalString(v)
	return res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalN__TypeKind2string(ctx context.Context, sel ast.SelectionSet, v string) graphql.Marshaler {
	res := graphql.MarshalString(v)
	if res == graphql.Null {
		if !graphql.HasFieldError(ctx, graphql.GetFieldContext(ctx)) {
			ec.Errorf(ctx, "the requested element is null which the schema does not allow")
		}
	}
	return res
}

func (ec *executionContext) unmarshalOBoolean2bool(ctx context.Context, v interface{}) (bool, error) {
	res, err := graphql.UnmarshalBoolean(v)
	return res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalOBoolean2bool(ctx context.Context, sel ast.SelectionSet, v bool) graphql.Marshaler {
	res := graphql.MarshalBoolean(v)
	return res
}

func (ec *executionContext) unmarshalOBoolean2ᚖbool(ctx context.Context, v interface{}) (*bool, error) {
	if v == nil {
		return nil, nil
	}
	res, err := graphql.UnmarshalBoolean(v)
	return &res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalOBoolean2ᚖbool(ctx context.Context, sel ast.SelectionSet, v *bool) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	res := graphql.MarshalBoolean(*v)
	return res
}

func (ec *executionContext) marshalOComment2ᚕᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐCommentᚄ(ctx context.Context, sel ast.SelectionSet, v []*models.Comment) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	ret := make(graphql.Array, len(v))
	var wg sync.WaitGroup
	isLen1 := len(v) == 1
	if !isLen1 {
		wg.Add(len(v))
	}
	for i := range v {
		i := i
		fc := &graphql.FieldContext{
			Index:  &i,
			Result: &v[i],
		}
		ctx := graphql.WithFieldContext(ctx, fc)
		f := func(i int) {
			defer func() {
				if r := recover(); r != nil {
					ec.Error(ctx, ec.Recover(ctx, r))
					ret = nil
				}
			}()
			if !isLen1 {
				defer wg.Done()
			}
			ret[i] = ec.marshalNComment2ᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋmodelsᚐComment(ctx, sel, v[i])
		}
		if isLen1 {
			f(i)
		} else {
			go f(i)
		}

	}
	wg.Wait()

	for _, e := range ret {
		if e == graphql.Null {
			return graphql.Null
		}
	}

	return ret
}

func (ec *executionContext) unmarshalOCommentFilter2ᚖgithubᚗcomᚋKilyakArataᚋgqlgenᚑtodosᚋgraphᚋmodelᚐCommentFilter(ctx context.Context, v interface{}) (*model.CommentFilter, error) {
	if v == nil {
		return nil, nil
	}
	res, err := ec.unmarshalInputCommentFilter(ctx, v)
	return &res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) unmarshalOID2ᚖstring(ctx context.Context, v interface{}) (*string, error) {
	if v == nil {
		return nil, nil
	}
	res, err := graphql.UnmarshalID(v)
	return &res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalOID2ᚖstring(ctx context.Context, sel ast.SelectionSet, v *string) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	res := graphql.MarshalID(*v)
	return res
}

func (ec *executionContext) unmarshalOInt2ᚖint(ctx context.Context, v interface{}) (*int, error) {
	if v == nil {
		return nil, nil
	}
	res, err := graphql.UnmarshalInt(v)
	return &res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalOInt2ᚖint(ctx context.Context, sel ast.SelectionSet, v *int) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	res := graphql.MarshalInt(*v)
	return res
}

func (ec *executionContext) unmarshalOString2ᚖstring(ctx context.Context, v interface{}) (*string, error) {
	if v == nil {
		return nil, nil
	}
	res, err := graphql.UnmarshalString(v)
	return &res, graphql.ErrorOnPath(ctx, err)
}

func (ec *executionContext) marshalOString2ᚖstring(ctx context.Context, sel ast.SelectionSet, v *string) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	res := graphql.MarshalString(*v)
	return res
}

func (ec *executionContext) marshalO__EnumValue2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐEnumValueᚄ(ctx context.Context, sel ast.SelectionSet, v []introspection.EnumValue) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	ret := make(graphql.Array, len(v))
	var wg sync.WaitGroup
	isLen1 := len(v) == 1
	if !isLen1 {
		wg.Add(len(v))
	}
	for i := range v {
		i := i
		fc := &graphql.FieldContext{
			Index:  &i,
			Result: &v[i],
		}
		ctx := graphql.WithFieldContext(ctx, fc)
		f := func(i int) {
			defer func() {
				if r := recover(); r != nil {
					ec.Error(ctx, ec.Recover(ctx, r))
					ret = nil
				}
			}()
			if !isLen1 {
				defer wg.Done()
			}
			ret[i] = ec.marshalN__EnumValue2githubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐEnumValue(ctx, sel, v[i])
		}
		if isLen1 {
			f(i)
		} else {
			go f(i)
		}

	}
	wg.Wait()

	for _, e := range ret {
		if e == graphql.Null {
			return graphql.Null
		}
	}

	return ret
}

func (ec *executionContext) marshalO__Field2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐFieldᚄ(ctx context.Context, sel ast.SelectionSet, v []introspection.Field) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	ret := make(graphql.Array, len(v))
	var wg sync.WaitGroup
	isLen1 := len(v) == 1
	if !isLen1 {
		wg.Add(len(v))
	}
	for i := range v {
		i := i
		fc := &graphql.FieldContext{
			Index:  &i,
			Result: &v[i],
		}
		ctx := graphql.WithFieldContext(ctx, fc)
		f := func(i int) {
			defer func() {
				if r := recover(); r != nil {
					ec.Error(ctx, ec.Recover(ctx, r))
					ret = nil
				}
			}()
			if !isLen1 {
				defer wg.Done()
			}
			ret[i] = ec.marshalN__Field2githubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐField(ctx, sel, v[i])
		}
		if isLen1 {
			f(i)
		} else {
			go f(i)
		}

	}
	wg.Wait()

	for _, e := range ret {
		if e == graphql.Null {
			return graphql.Null
		}
	}

	return ret
}

func (ec *executionContext) marshalO__InputValue2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐInputValueᚄ(ctx context.Context, sel ast.SelectionSet, v []introspection.InputValue) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	ret := make(graphql.Array, len(v))
	var wg sync.WaitGroup
	isLen1 := len(v) == 1
	if !isLen1 {
		wg.Add(len(v))
	}
	for i := range v {
		i := i
		fc := &graphql.FieldContext{
			Index:  &i,
			Result: &v[i],
		}
		ctx := graphql.WithFieldContext(ctx, fc)
		f := func(i int) {
			defer func() {
				if r := recover(); r != nil {
					ec.Error(ctx, ec.Recover(ctx, r))
					ret = nil
				}
			}()
			if !isLen1 {
				defer wg.Done()
			}
			ret[i] = ec.marshalN__InputValue2githubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐInputValue(ctx, sel, v[i])
		}
		if isLen1 {
			f(i)
		} else {
			go f(i)
		}

	}
	wg.Wait()

	for _, e := range ret {
		if e == graphql.Null {
			return graphql.Null
		}
	}

	return ret
}

func (ec *executionContext) marshalO__Schema2ᚖgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐSchema(ctx context.Context, sel ast.SelectionSet, v *introspection.Schema) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	return ec.___Schema(ctx, sel, v)
}

func (ec *executionContext) marshalO__Type2ᚕgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐTypeᚄ(ctx context.Context, sel ast.SelectionSet, v []introspection.Type) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	ret := make(graphql.Array, len(v))
	var wg sync.WaitGroup
	isLen1 := len(v) == 1
	if !isLen1 {
		wg.Add(len(v))
	}
	for i := range v {
		i := i
		fc := &graphql.FieldContext{
			Index:  &i,
			Result: &v[i],
		}
		ctx := graphql.WithFieldContext(ctx, fc)
		f := func(i int) {
			defer func() {
				if r := recover(); r != nil {
					ec.Error(ctx, ec.Recover(ctx, r))
					ret = nil
				}
			}()
			if !isLen1 {
				defer wg.Done()
			}
			ret[i] = ec.marshalN__Type2githubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐType(ctx, sel, v[i])
		}
		if isLen1 {
			f(i)
		} else {
			go f(i)
		}

	}
	wg.Wait()

	for _, e := range ret {
		if e == graphql.Null {
			return graphql.Null
		}
	}

	return ret
}

func (ec *executionContext) marshalO__Type2ᚖgithubᚗcomᚋ99designsᚋgqlgenᚋgraphqlᚋintrospectionᚐType(ctx context.Context, sel ast.SelectionSet, v *introspection.Type) graphql.Marshaler {
	if v == nil {
		return graphql.Null
	}
	return ec.___Type(ctx, sel, v)
}
