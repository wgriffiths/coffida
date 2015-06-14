package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/manyminds/api2go"
	"github.com/manyminds/api2go/jsonapi"
)

type Product struct {
	ID          string
	Title       string  
	Description string  
	Price       int32
	Active      bool
}

// GetID to satisfy jsonapi.MarshalIdentifier interface
func (p Product) GetID() string {
	return p.ID
}

// SetID to satisfy jsonapi.UnmarshalIdentifier interface
func (p *Product) SetID(id string) error {
	p.ID = id
	return nil
}

// GetReferences to satisfy the jsonapi.MarshalReferences interface
func (p Product) GetReferences() []jsonapi.Reference {
	return []jsonapi.Reference{}
}

// GetReferencedIDs to satisfy the jsonapi.MarshalLinkedRelations interface
func (p Product) GetReferencedIDs() []jsonapi.ReferenceID {
	result := []jsonapi.ReferenceID{}
	return result
}

// GetReferencedStructs to satisfy the jsonapi.MarhsalIncludedRelations interface
func (p Product) GetReferencedStructs() []jsonapi.MarshalIdentifier {
	result := []jsonapi.MarshalIdentifier{}
	return result
}

// SetToManyReferenceIDs sets the sweets reference IDs and satisfies the jsonapi.UnmarshalToManyRelations interface
func (p *Product) SetToManyReferenceIDs(name string, IDs []string) error {
	return errors.New("There is no to-many relationship with the name " + name)
}

// AddToManyIDs adds some new sweets that a users loves so much
func (p *Product) AddToManyIDs(name string, IDs []string) error {
	return errors.New("There is no to-many relationship with the name " + name)
}

// DeleteToManyIDs removes some sweets from a users because they made him very sick
func (p *Product) DeleteToManyIDs(name string, IDs []string) error {
	return errors.New("There is no to-many relationship with the name " + name)
}

// the user resource holds all users in the array
type productResource struct {
	products map[string]Product
	idCount  int
}

// FindAll to satisfy api2go data source interface
func (p *productResource) FindAll(r api2go.Request) (interface{}, error) {
	var products []Product

	for _, value := range p.products {
		products = append(products, value)
	}

	return products, nil
}

func (p *productResource) PaginatedFindAll(r api2go.Request) (interface{}, uint, error) {
	var (
		products                    []Product
		number, size, offset, limit string
		keys                        []int
	)

	for k := range p.products {
		i, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			return nil, 0, err
		}

		keys = append(keys, int(i))
	}
	sort.Ints(keys)

	numberQuery, ok := r.QueryParams["page[number]"]
	if ok {
		number = numberQuery[0]
	}
	sizeQuery, ok := r.QueryParams["page[size]"]
	if ok {
		size = sizeQuery[0]
	}
	offsetQuery, ok := r.QueryParams["page[offset]"]
	if ok {
		offset = offsetQuery[0]
	}
	limitQuery, ok := r.QueryParams["page[limit]"]
	if ok {
		limit = limitQuery[0]
	}

	if size != "" {
		sizeI, err := strconv.ParseUint(size, 10, 64)
		if err != nil {
			return nil, 0, err
		}

		numberI, err := strconv.ParseUint(number, 10, 64)
		if err != nil {
			return nil, 0, err
		}

		start := sizeI * (numberI - 1)
		for i := start; i < start+sizeI; i++ {
			if i >= uint64(len(p.products)) {
				break
			}
			products = append(products, p.products[strconv.FormatInt(int64(keys[i]), 10)])
		}
	} else {
		limitI, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			return nil, 0, err
		}

		offsetI, err := strconv.ParseUint(offset, 10, 64)
		if err != nil {
			return nil, 0, err
		}

		for i := offsetI; i < offsetI+limitI; i++ {
			if i >= uint64(len(p.products)) {
				break
			}
			products = append(products, p.products[strconv.FormatInt(int64(keys[i]), 10)])
		}
	}

	return products, uint(len(p.products)), nil
}

// FindOne to satisfy `api2go.DataSource` interface
// this method should return the user with the given ID, otherwise an error
func (p *productResource) FindOne(ID string, r api2go.Request) (interface{}, error) {
	if product, ok := p.products[ID]; ok {
		return product, nil
	}

	return nil, api2go.NewHTTPError(errors.New("Not Found"), "Not Found", http.StatusNotFound)
}

// Create method to satisfy `api2go.DataSource` interface
func (s *productResource) Create(obj interface{}, r api2go.Request) (string, error) {
	product, ok := obj.(Product)
	fmt.Printf("%v\n", r)
	fmt.Printf("%v\n", obj)
	fmt.Printf("%v\n", product)
	if !ok {
		return "", api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	if _, ok := s.products[product.GetID()]; ok {
		return "", api2go.NewHTTPError(errors.New("Product exists"), "Product exists", http.StatusConflict)
	}

	s.idCount++
	id := fmt.Sprintf("%d", s.idCount)
	product.SetID(id)

	s.products[id] = product

	return id, nil
}

// Delete to satisfy `api2go.DataSource` interface
func (s *productResource) Delete(id string, r api2go.Request) error {
	obj, err := s.FindOne(id, api2go.Request{})
	if err != nil {
		return err
	}

	product, ok := obj.(Product)
	if !ok {
		return errors.New("Invalid instance given")
	}

	delete(s.products, product.GetID())

	return nil
}

//Update stores all changes on the user
func (s *productResource) Update(obj interface{}, r api2go.Request) error {
	product, ok := obj.(Product)
	if !ok {
		return api2go.NewHTTPError(errors.New("Invalid instance given"), "Invalid instance given", http.StatusBadRequest)
	}

	s.products[product.GetID()] = product

	return nil
}

// PrettyJSONContentMarshaler for JSON in a human readable format
type PrettyJSONContentMarshaler struct{}

// Marshal marshals to pretty JSON
func (m PrettyJSONContentMarshaler) Marshal(i interface{}) ([]byte, error) {
	return json.MarshalIndent(i, "", "    ")
}

// Unmarshal the JSON
func (m PrettyJSONContentMarshaler) Unmarshal(data []byte, i interface{}) error {
	return json.Unmarshal(data, i)
}

func main() {
	marshalers := map[string]api2go.ContentMarshaler{
		"application/vnd.api+json": PrettyJSONContentMarshaler{},
	}

	api := api2go.NewAPIWithMarshalers("v0", "http://localhost:31415", marshalers)
	products := make(map[string]Product)
	api.AddResource(Product{}, &productResource{products: products})

	fmt.Println("Listening on :31415")
	handler := api.Handler().(*httprouter.Router)
	http.ListenAndServe(":31415", handler)
}
